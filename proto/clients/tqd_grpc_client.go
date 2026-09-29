// common/protobuf/clients/tqd_grpc_client.go
package clients

import (
	"context"
	"errors"

	tqdpb "pb/types/tqd"

	"google.golang.org/grpc"
)

type TQDGrpcClient struct {
	Client              tqdpb.LocationServiceClient
	ParcelClient        tqdpb.ParcelServiceClient
	SeoProjectionClient tqdpb.TqdSeoProjectionServiceClient
}

func BindTQDClient(conn grpc.ClientConnInterface) *TQDGrpcClient {
	return &TQDGrpcClient{
		Client:              tqdpb.NewLocationServiceClient(conn),
		ParcelClient:        tqdpb.NewParcelServiceClient(conn),
		SeoProjectionClient: tqdpb.NewTqdSeoProjectionServiceClient(conn),
	}
}

// func NewTQDClient(ep config.ServiceEndpoint) (*TQDGrpcClient, error) {
//     if err != nil {
//         return nil, err
//     }
//     return &TQDGrpcClient{
//         Client:         tqdpb.NewLocationServiceClient(base.Conn()),
//     }, nil
// }

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
