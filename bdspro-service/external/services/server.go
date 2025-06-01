package service

// import (
// 	"bdspro/internal/domain"
// 	"bdspro/internal/usecases"
// 	"context"
// 	pb_bdspro "pb/types/bdspro"

// 	"github.com/jinzhu/copier"
// )

// type AssetCostServer struct {
// 	pb_bdspro.UnimplementedAssetCostServiceServer
// 	UC usecases.AssetCostUsecase
// }

// func NewAssetCostServer(AssetCostUc usecases.AssetCostUsecase) *AssetCostServer {
// 	return &AssetCostServer{}
// }

// func (s AssetCostServer) Get(c context.Context, dto *pb_bdspro.AssetCostDTO) (*pb_bdspro.AssetCostListDTO, error) {
// 	result, err := s.UC.GetData(c, dto.Id)

// 	pbdatas := []*pb_bdspro.AssetCostDTO{}
// 	for _, r := range result {
// 		pbdata := &pb_bdspro.AssetCostDTO{}
// 		copier.Copy(r, &pbdata)
// 		pbdatas = append(pbdatas, pbdata)
// 	}

// 	return &pb_bdspro.AssetCostListDTO{
// 		Data: pbdatas,
// 	}, err
// }
// func (s AssetCostServer) Detail(c context.Context, dto *pb_bdspro.AssetCostDTO) (*pb_bdspro.AssetCostDTO, error) {
// 	r, err := s.UC.GetByID(c, dto.Id)

// 	pbdata := &pb_bdspro.AssetCostDTO{}
// 	copier.Copy(r, &pbdata)

// 	return pbdata, err
// }
// func (s AssetCostServer) Create(c context.Context, dto *pb_bdspro.AssetCostDTO) (*pb_bdspro.AssetCostDTO, error) {
// 	d := &domain.AssetCost{}
// 	copier.Copy(&d, &dto)

// 	r, err := s.UC.Create(c, d)

// 	copier.Copy(&dto, &r)
// 	return dto, err
// }
// func (s AssetCostServer) Update(c context.Context, dto *pb_bdspro.AssetCostDTO) (*pb_bdspro.AssetCostDTO, error) {
// 	d := &domain.AssetCost{}
// 	copier.Copy(&d, &dto)

// 	r, err := s.UC.Update(c, dto.Id, d)

// 	copier.Copy(&dto, &r)
// 	return dto, err
// }
// func (s AssetCostServer) Delete(c context.Context, dto *pb_bdspro.AssetCostDTO) (*pb_bdspro.AssetCostDTO, error) {
// 	err := s.UC.Delete(c, dto.Id)
// 	return dto, err
// }
