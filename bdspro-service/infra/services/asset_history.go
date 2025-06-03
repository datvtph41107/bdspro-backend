package service

import (
	"bdspro/internal/usecases"
	pb_bdspro "pb/types/bdspro"
)

type AssetSplitHistoryServer struct {
	pb_bdspro.UnimplementedAssetSplitHistoryServiceServer
	UC usecases.AssetHistoryUsecase
}

func NewAssetSplitHistoryServer(AssetSplitHistoryUc usecases.AssetHistoryUsecase) *AssetSplitHistoryServer {
	return &AssetSplitHistoryServer{
		UC: AssetSplitHistoryUc,
	}
}

// func (s AssetSplitHistoryServer) Get(c context.Context, dto *pb_bdspro.AssetSplitHistoryDTO) (*pb_bdspro.AssetSplitHistoryListDTO, error) {
// 	result, err := s.UC.GetData(c, dto.Id)

// 	pbdatas := []*pb_bdspro.AssetSplitHistoryDTO{}
// 	for _, r := range result {
// 		pbdata := &pb_bdspro.AssetSplitHistoryDTO{}
// 		copier.Copy(&pbdata, &r)
// 		pbdatas = append(pbdatas, pbdata)
// 	}

// 	return &pb_bdspro.AssetSplitHistoryListDTO{
// 		Data: pbdatas,
// 	}, err
// }
// func (s AssetSplitHistoryServer) Detail(c context.Context, dto *pb_bdspro.AssetSplitHistoryDTO) (*pb_bdspro.AssetSplitHistoryDTO, error) {
// 	r, err := s.UC.GetByID(c, dto.Id)

// 	pbdata := &pb_bdspro.AssetSplitHistoryDTO{}
// 	copier.Copy(&pbdata, &r)

// 	return pbdata, err
// }
// func (s AssetSplitHistoryServer) Create(c context.Context, dto *pb_bdspro.AssetSplitHistoryDTO) (*pb_bdspro.AssetSplitHistoryDTO, error) {
// 	d := &domain.AssetHistory{}
// 	copier.Copy(&d, &dto)

// 	r, err := s.UC.Create(c, d)

// 	copier.Copy(&dto, &r)
// 	return dto, err
// }
// func (s AssetSplitHistoryServer) Update(c context.Context, dto *pb_bdspro.AssetSplitHistoryDTO) (*pb_bdspro.AssetSplitHistoryDTO, error) {
// 	d := &domain.AssetHistory{}
// 	copier.Copy(&d, &dto)

// 	r, err := s.UC.Update(c, dto.Id, d)

// 	copier.Copy(&dto, &r)
// 	return dto, err
// }
// func (s AssetSplitHistoryServer) Delete(c context.Context, dto *pb_bdspro.AssetSplitHistoryDTO) (*pb_bdspro.AssetSplitHistoryDTO, error) {
// 	err := s.UC.Delete(c, dto.Id)
// 	return dto, err
// }
