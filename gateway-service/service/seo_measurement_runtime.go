package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"common"
	_middleware "common/middleware"
	crmpb "pb/types/crm"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

const seoM8BodyLimit = 8 << 20

func normalizeJSONNumbers(value interface{}) interface{} {
	switch typed := value.(type) {
	case json.Number:
		if integer, err := typed.Int64(); err == nil {
			return float64(integer)
		}
		if decimal, err := typed.Float64(); err == nil {
			return decimal
		}
		return typed.String()
	case map[string]interface{}:
		out := make(map[string]interface{}, len(typed))
		for key, item := range typed {
			out[key] = normalizeJSONNumbers(item)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(typed))
		for index, item := range typed {
			out[index] = normalizeJSONNumbers(item)
		}
		return out
	default:
		return value
	}
}

func DialCRMGRPC(env string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	port := viper.GetString("service.port.crm")
	host := common.GetPrefixProtobufUrl(env, "crm")
	return grpc.NewClient(fmt.Sprintf("%s:%s", host, port), opts...)
}

func seoM8Context(request *http.Request) context.Context {
	md := _middleware.InjectGrpcMetadataContextMiddleware(request.Context(), request)
	if existing, ok := metadata.FromOutgoingContext(request.Context()); ok {
		md = metadata.Join(existing, md)
	}
	return metadata.NewOutgoingContext(request.Context(), md)
}

func seoM8RequestFromBody(writer http.ResponseWriter, request *http.Request) (*structpb.Struct, error) {
	reader := http.MaxBytesReader(writer, request.Body, seoM8BodyLimit)
	defer reader.Close()
	var payload map[string]interface{}
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	if err := decoder.Decode(&payload); err != nil {
		if err == io.EOF {
			payload = map[string]interface{}{}
		} else {
			return nil, err
		}
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("request body must contain one JSON object")
		}
		return nil, err
	}
	return structpb.NewStruct(normalizeJSONNumbers(payload).(map[string]interface{}))
}

func seoM8QueryStruct(request *http.Request) (*structpb.Struct, error) {
	query := request.URL.Query()
	payload := map[string]interface{}{}
	for _, key := range []string{"source", "status", "grade", "canonicalUrl", "range"} {
		if value := strings.TrimSpace(query.Get(key)); value != "" {
			payload[key] = value
		}
	}
	for _, key := range []string{"page", "size", "limit"} {
		if value := strings.TrimSpace(query.Get(key)); value != "" {
			parsed, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("%s must be an integer", key)
			}
			payload[key] = float64(parsed)
		}
	}
	return structpb.NewStruct(payload)
}

func writeSeoM8HTTP(writer http.ResponseWriter, response *structpb.Struct, err error) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("Cache-Control", "no-store")
	if err != nil {
		writer.WriteHeader(discoveryHTTPStatus(err))
		message := err.Error()
		if grpcStatus, ok := status.FromError(err); ok {
			message = grpcStatus.Message()
		}
		_ = json.NewEncoder(writer).Encode(map[string]interface{}{"code": status.Code(err).String(), "message": message})
		return
	}
	if response == nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(writer).Encode(map[string]interface{}{"code": "EMPTY_RESPONSE", "message": "CRM SEO-M8 returned an empty response"})
		return
	}
	payload, marshalErr := protojson.Marshal(response)
	if marshalErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(writer).Encode(map[string]interface{}{"code": "ENCODING_FAILED", "message": marshalErr.Error()})
		return
	}
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(payload)
}

func seoM8Call(writer http.ResponseWriter, request *http.Request, timeout time.Duration, call func(context.Context) (*structpb.Struct, error)) {
	ctx, cancel := context.WithTimeout(seoM8Context(request), timeout)
	defer cancel()
	response, err := call(ctx)
	writeSeoM8HTTP(writer, response, err)
}

// RegisterSeoMeasurementRuntimeHTTP exposes the M8 control plane through the
// existing authenticated Gateway /v2 catch-all while keeping CRM transport gRPC-only.
func RegisterSeoMeasurementRuntimeHTTP(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) error {
	conn, err := DialCRMGRPC(env, opts)
	if err != nil {
		return err
	}
	client := crmpb.NewSeoMeasurementRuntimeServiceClient(conn)

	if err := mux.HandlePath(http.MethodPost, "/v2/crm/admin/seo/measurement/m8/import/{source}", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		input, parseErr := seoM8RequestFromBody(w, r)
		if parseErr != nil {
			writeSeoM8HTTP(w, nil, status.Error(codes.InvalidArgument, parseErr.Error()))
			return
		}
		if input.Fields == nil {
			input.Fields = map[string]*structpb.Value{}
		}
		input.Fields["source"] = structpb.NewStringValue(params["source"])
		seoM8Call(w, r, 45*time.Second, func(callCtx context.Context) (*structpb.Struct, error) {
			return client.ImportMeasurement(callCtx, input)
		})
	}); err != nil {
		return err
	}

	if err := mux.HandlePath(http.MethodGet, "/v2/crm/admin/seo/measurement/m8/sources", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		seoM8Call(w, r, 10*time.Second, func(callCtx context.Context) (*structpb.Struct, error) {
			return client.GetSourceHealth(callCtx, &structpb.Struct{})
		})
	}); err != nil {
		return err
	}

	if err := mux.HandlePath(http.MethodGet, "/v2/crm/admin/seo/measurement/m8/runs", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		input, parseErr := seoM8QueryStruct(r)
		if parseErr != nil {
			writeSeoM8HTTP(w, nil, status.Error(codes.InvalidArgument, parseErr.Error()))
			return
		}
		seoM8Call(w, r, 10*time.Second, func(callCtx context.Context) (*structpb.Struct, error) { return client.ListImportRuns(callCtx, input) })
	}); err != nil {
		return err
	}

	if err := mux.HandlePath(http.MethodPost, "/v2/crm/admin/seo/measurement/m8/score", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		input, parseErr := seoM8RequestFromBody(w, r)
		if parseErr != nil {
			writeSeoM8HTTP(w, nil, status.Error(codes.InvalidArgument, parseErr.Error()))
			return
		}
		seoM8Call(w, r, 45*time.Second, func(callCtx context.Context) (*structpb.Struct, error) { return client.RunScoring(callCtx, input) })
	}); err != nil {
		return err
	}

	if err := mux.HandlePath(http.MethodGet, "/v2/crm/admin/seo/measurement/m8/scores", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		input, parseErr := seoM8QueryStruct(r)
		if parseErr != nil {
			writeSeoM8HTTP(w, nil, status.Error(codes.InvalidArgument, parseErr.Error()))
			return
		}
		seoM8Call(w, r, 10*time.Second, func(callCtx context.Context) (*structpb.Struct, error) { return client.ListPageScores(callCtx, input) })
	}); err != nil {
		return err
	}

	_ = ctx
	return nil
}
