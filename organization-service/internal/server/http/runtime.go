package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/interface/handler"
	"github.com/tuandq2112/organization-service/internal/server/http/middlewares"
	organizationpb "pb/types/organization"
)

type HttpServer struct {
	port                string
	organizationHandler *handler.OrganizationHandler
	logger              *flogging.FabricLogger
}

func NewHttpServer(organizationHandler *handler.OrganizationHandler) *HttpServer {
	return &HttpServer{
		port:                appconfig.HTTP_PORT,
		organizationHandler: organizationHandler,
		logger:              flogging.MustGetLogger("http_service"),
	}
}


func (h *HttpServer) Run() {
	ctx := context.Background()
	var opts []runtime.ServeMuxOption
	opts = append(opts, runtime.WithErrorHandler(CustomErrorHandler))
	mux := runtime.NewServeMux(opts...)
	err := organizationpb.RegisterOrganizationServiceHandlerServer(ctx, mux, h.organizationHandler)
	if err != nil {
		h.logger.Errorf("Failed to register organization service: %v", err)
		return
	}

	swaggerMux := http.NewServeMux()
	swaggerMux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger/" {
			http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
			return
		}
		if strings.HasSuffix(r.URL.Path, ".json") {
			relPath := strings.TrimPrefix(r.URL.Path, "/swagger/")
			filePath := "proto/docs/" + relPath
			http.ServeFile(w, r, filePath)
			return
		}
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/root.swagger.json"),
		).ServeHTTP(w, r)
	})

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/swagger/") {
			swaggerMux.ServeHTTP(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})

	handler = middlewares.JWTMiddleware(handler)
	handler = middlewares.LoggingMiddleware(handler)
	handler = middlewares.EnableCORS(handler)

	server := &http.Server{
		Addr:    ":" + h.port,
		Handler: handler,
	}
	h.logger.Infof("Server is running on port %s", h.port)
	err = server.ListenAndServe()
	if err != nil {
		h.logger.Errorf("Failed to listen and serve: %v", err)
		return
	}
}
