// common/protobuf/clients/tqd_grpc_client.go
package clients

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"common"
	_middleware "common/middleware"
	tqdpb "pb/types/tqd"

	"google.golang.org/grpc"
)

type TQDGrpcClient struct {
	Client tqdpb.LocationServiceClient
	conn   *grpc.ClientConn
}

var (
	tqdClientInstance *TQDGrpcClient
	tqdOnce           sync.Once
)

// type TQDGrpcClient struct {
//     *BaseGrpcClient
//     Client tqdpb.LocationServiceClient
// }

// func NewTQDClient(ep config.ServiceEndpoint) (*TQDGrpcClient, error) {
//     base, err := NewBaseGrpcClient(ep)
//     if err != nil {
//         return nil, err
//     }
//     return &TQDGrpcClient{
//         BaseGrpcClient: base,
//         Client:         tqdpb.NewLocationServiceClient(base.Conn()),
//     }, nil
// }

// GetTQDClient - Singleton pattern giống user client của bạn
func GetTQDClient(envRuntime string) *TQDGrpcClient {
	tqdOnce.Do(func() {
		client := &TQDGrpcClient{}
		go client.initConnection(envRuntime)
		tqdClientInstance = client
	})
	return tqdClientInstance
}

// NewTQDClient - Tạo instance mới (không singleton)
func NewTQDClient(envRuntime string) *TQDGrpcClient {
	client := &TQDGrpcClient{}

	go client.initConnection(envRuntime)
	return client
}

func (c *TQDGrpcClient) initConnection(envRuntime string) {
	name := "tqd"
	log.Printf("User URL = %s", common.GetPrefixProtobufUrl(envRuntime, "user"))
	log.Printf("TQD URL = %s", common.GetPrefixProtobufUrl(envRuntime, "tqd"))
	url := common.GetPrefixProtobufUrl(envRuntime, name)
	log.Printf("Connecting to TQD service at %s:8219", url)
	conn, err := grpc.Dial("localhost:8219",
		grpc.WithInsecure(),
		grpc.WithBackoffMaxDelay(5*time.Second),
		grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
	)
	if err != nil {
		log.Fatalf("failed to dial TQD service: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Monitor connection như user client
	go _middleware.MonitorConnection(ctx, name, conn)

	c.conn = conn
	c.Client = tqdpb.NewLocationServiceClient(conn)
}

// Close - Đóng connection
func (c *TQDGrpcClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetNearestLocation - Lấy province/ward gần nhất từ tọa độ
func (c *TQDGrpcClient) GetNearestLocation(ctx context.Context, lat, lng float64) (*tqdpb.LocationResponse, error) {
	if c.Client == nil {
		return nil, errors.New("TQD client not available")
	}

	req := &tqdpb.GetNearestLocationRequest{
		Latitude:      lat,
		Longitude:     lng,
		MaxDistanceKm: nil, // dùng default 100km
		Limit:         nil, // dùng default 1
	}

	return c.Client.GetNearestLocation(ctx, req)
}

// BatchGetNearestLocations - Lấy nhiều location cùng lúc
func (c *TQDGrpcClient) BatchGetNearestLocations(ctx context.Context, coordinates []*tqdpb.Coordinate) (*tqdpb.BatchLocationResponse, error) {
	if c.Client == nil {
		return nil, errors.New("TQD client not available")
	}

	req := &tqdpb.BatchGetNearestLocationsRequest{
		Coordinates: coordinates,
	}

	return c.Client.BatchGetNearestLocations(ctx, req)
}
