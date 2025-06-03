package handler

import (
	"bdspro/internal/usecases"
	pb_generic "pb/types/generic"
)

type GenericGrpcServer struct {
	pb_generic.UnimplementedGenericServiceServer
	// AssetCostHandler AssetCostHandler
}

func NewGenericService(asset *usecases.UAssetUsecase,

// AssetCostHandler AssetCostHandler,
) *GenericGrpcServer {
	return &GenericGrpcServer{
		// AssetCostHandler: AssetCostHandler,
	}
}

// func (s GenericGrpcServer) HandleRequest(c context.Context, req *pb_generic.RequestBody) (*pb_generic.ResponseBody, error) {
// 	log.Println("HANDLE_REQUEST")
// 	switch req.Path {
// 	case "/":
// 		result, err := s.AssetCostHandler.RegisterHandler(c, req)
// 		return &pb_generic.ResponseBody{Body: result}, err
// 	case "/2":
// 		{

// 		}
// 	}
// 	// if result, err := s.AssetCostHandler.RegisterHandler(c, req); result != nil || err != nil {
// 	// 	var body *anypb.Any
// 	// 	if result != nil {
// 	// 		an, err := anypb.New(&structpb.Struct{
// 	// 			Fields: map[string]*structpb.Value{
// 	// 				"content": structpb.NewStringValue("result"),
// 	// 			},
// 	// 		})
// 	// 		if err != nil {
// 	// 			return nil, err
// 	// 		}
// 	// 		body = an
// 	// 	}
// 	// 	return &pb_generic.ResponseBody{
// 	// 		Body: body,
// 	// 	}, err
// 	// }
// 	return nil, status.Errorf(codes.Unimplemented, "method HandleRequest not implemented")
// }
