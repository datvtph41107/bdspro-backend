package handler

// type BdsproGrpcServer struct {
// 	pb_bdspro.UnimplementedBdsproServiceServer
// 	AssetUsecase *usecases.UAssetUsecase
// 	AssetCostUc  *usecases.AssetCostUsecase
// }

// func NewGrpcService(asset *usecases.UAssetUsecase,
// 	AssetCostUc *usecases.AssetCostUsecase,
// ) *BdsproGrpcServer {
// 	return &BdsproGrpcServer{
// 		AssetUsecase: asset,
// 		AssetCostUc:  AssetCostUc,
// 	}
// }

// func (s *BdsproGrpcServer) GetAssets(c context.Context, dto *pb_bdspro.Empty) (*pb_bdspro.AssetList, error) {
// 	profileId := _utils.GetProfileIdWithContext(c)
// 	result, err := s.AssetUsecase.GetOwnerAssets(c, profileId, 20, 0)
// 	assets := []*pb_bdspro.Asset{}
// 	for _, asset := range result {
// 		pbdata := &pb_bdspro.Asset{}
// 		copier.Copy(asset, &pbdata)
// 		assets = append(assets, pbdata)
// 	}
// 	return &pb_bdspro.AssetList{
// 		Data: assets,
// 	}, err
// }

// // func (BdsproGrpcServer) FilterAssets(context.Context, *pb_bdspro.FilterAssetsRequest) (*pb_bdspro.AssetList, error) {
// // 	return nil, status.Errorf(codes.Unimplemented, "method FilterAssets not implemented")
// // }
// // func ConvertToProtoAsset(asset *domain.Asset) *pb_bdspro.Asset {
// // 	if asset == nil {
// // 		return nil
// // 	}

// // 	return &pb_bdspro.Asset{
// // 		Id:               asset.ID,
// // 		Name:             asset.Name,
// // 		Description:      asset.Description,
// // 		Area:             asset.Area,
// // 		Address:          asset.Address,
// // 		WardId:           asset.WardId,
// // 		DistrictId:       asset.DistrictID,
// // 		ProvinceId:       asset.ProvinceId,
// // 		Price:            asset.Price,
// // 		BuyedDate:        _utils.TimeToTimestamp(asset.BuyedDate),
// // 		DocStatus:        asset.DocStatus,
// // 		ProductId:        asset.ProductID,
// // 		OwnerUserId:      asset.OwnerUserID,
// // 		OwnerOrgId:       asset.OwnerOrgID,
// // 		ParentAssetId:    asset.ParentAssetID,
// // 		SplitMergeStatus: asset.SplitMergeStatus,
// // 		PurchasePrice:    asset.PurchasePrice,
// // 		PurchaseDate:     _utils.TimeToTimestamp(asset.PurchaseDate),
// // 		LegalStatus:      asset.LegalStatus,
// // 		Archived:         asset.Archived,
// // 		RentStatus:       uint32(asset.RentStatus),
// // 		// Status:           asset.Status,
// // 		// Type:             asset.Type,
// // 	}
// // }
// // func (s BdsproGrpcServer) BulkAction(ctx context.Context, req *pb_bdspro.BulkActionRequest) (*pb_bdspro.Empty, error) {
// // 	// var ids []uint64
// // 	// for _, id := range req.AssetId {
// // 	// 	ids = append(ids, id)
// // 	// }

// // 	// if err := s.AssetUsecase.BulkAction(ctx, req.Action, ids); err != nil {
// // 	// 	return nil, status.Errorf(codes.Internal, "bulk action failed: %v", err)
// // 	// }

// // 	return &pb_bdspro.Empty{}, nil
// // }

// // func (s BdsproGrpcServer) CreateAsset(ctx context.Context, req *pb_bdspro.CreateAssetRequest) (*pb_bdspro.Asset, error) {
// // 	asset := &domain.Asset{}

// // 	if err := _utils.MapWithValidation(req, asset); err != nil {
// // 		return nil, _utils.ErrorField(err, &pb_spec.FieldErrors{
// // 			Errors: err.Errors,
// // 		})
// // 	}

// // 	result, err := s.AssetUsecase.CreateAsset(ctx, asset)
// // 	// if err != nil {
// // 	// 	return nil, status.Errorf(codes.Internal, "failed to create asset: %v", err)
// // 	// }

// // 	return ConvertToProtoAsset(result), err
// // }

// // func (s *BdsproGrpcServer) CreateAssetFromProduct(ctx context.Context, req *pb_bdspro.CreateAssetRequestFromProduct) (*pb_bdspro.Asset, error) {
// // 	asset := &domain.Asset{}

// // 	if err := _utils.MapWithValidation(req, asset); err != nil {
// // 		return nil, _utils.ErrorField(err, &pb_spec.FieldErrors{
// // 			Errors: err.Errors,
// // 		})
// // 	}

// // 	result, err := s.AssetUsecase.CreateAssetFromProduct(ctx, req.ProductId, asset) //CreateAssetFromProduct
// // 	return ConvertToProtoAsset(result), err
// // }

// // func (s *BdsproGrpcServer) ConvertAssetFromProduct(ctx context.Context, req *pb_bdspro.CreateAssetRequest) (*pb_bdspro.Asset, error) {
// // 	asset := &domain.Asset{}

// // 	if err := _utils.MapWithValidation(req, asset); err != nil {
// // 		return nil, _utils.ErrorField(err, &pb_spec.FieldErrors{
// // 			Errors: err.Errors,
// // 		})
// // 	}

// // 	result, err := s.AssetUsecase.ConvertProductToAsset(ctx, *req.ProductId, asset) //CreateAssetFromProduct
// // 	return ConvertToProtoAsset(result), err

// // }

// // // func (s *GrpcServer) ConvertProductToAsset(ctx context.Context, productID uint64, req *pb_bdspro.CreateAssetRequest) (*pb_bdspro.Asset, error) {
// // // 	product, err := s.productRepo.GetByID(productID)
// // // 	if err != nil {
// // // 		return nil, fmt.Errorf("failed to get product: %w", err)
// // // 	}

// // // 	// Điền mặc định thông tin từ product nếu không có trong req
// // // 	if req.Name == "" {
// // // 		req.Name = product.Name
// // // 	}
// // // 	if req.Area == 0 {
// // // 		req.Area = product.Area
// // // 	}
// // // 	if req.Address == "" {
// // // 		req.Address = product.Address
// // // 	}

// // // 	req.ProductId = &productID
// // // 	return s.saveAsset(ctx, req)
// // // }

// // func (s BdsproGrpcServer) UpdateAsset(ctx context.Context, req *pb_bdspro.UpdateAssetRequest) (*pb_bdspro.Asset, error) {
// // 	// Chuyển đổi google.protobuf.Timestamp sang time.Time nếu cần
// // 	var purchaseDate time.Time
// // 	if req.PurchaseDate != nil {
// // 		purchaseDate = req.PurchaseDate.AsTime()
// // 	}

// // 	asset := &domain.Asset{
// // 		// ID:          req.Id,
// // 		Name:        req.Name,
// // 		Description: req.Description,
// // 		Area:        req.Area,
// // 		Address:     req.Address,
// // 		// OwnerID:          req.OwnerId,
// // 		// Status:           req.Status,
// // 		// Type:             req.Type,
// // 		WardId:           req.WardId,
// // 		DistrictId:       req.DistrictId,
// // 		ProvinceId:       req.ProvinceId,
// // 		Price:            req.Price,
// // 		BuyedDate:        _utils.TimestampToTime(req.BuyedDate),
// // 		DocStatus:        req.DocStatus,
// // 		ProductID:        req.ProductId,
// // 		OwnerUserID:      req.OwnerUserId,
// // 		OwnerOrgID:       req.OwnerOrgId,
// // 		ParentAssetID:    req.ParentAssetId,
// // 		SplitMergeStatus: req.SplitMergeStatus,
// // 		PurchasePrice:    req.PurchasePrice,
// // 		PurchaseDate:     &purchaseDate,
// // 		LegalStatus:      req.LegalStatus,
// // 		// Archived:         req.Archived,
// // 		RentStatus: enums.AssetStatus(req.RentStatus),
// // 	}

// // 	// Gọi service cập nhật tài sản
// // 	err := s.AssetUsecase.UpdateAsset(ctx, req.Id, asset)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	// Chuyển đổi lại về protobuf response
// // 	return ConvertToProtoAsset(asset), nil
// // }

// // func (s BdsproGrpcServer) DeleteAsset(ctx context.Context, req *pb_bdspro.DeleteAssetRequest) (*pb_bdspro.Empty, error) {
// // 	err := s.AssetUsecase.DeleteAsset(ctx, req.GetId())
// // 	if err != nil {
// // 		return nil, status.Errorf(codes.Internal, "failed to delete asset: %v", err)
// // 	}
// // 	return &pb_bdspro.Empty{}, nil
// // }

// // func (s BdsproGrpcServer) GetAssetDetails(ctx context.Context, req *pb_bdspro.GetAssetRequest) (*pb_bdspro.AssetDetails, error) {
// // 	asset, err := s.AssetUsecase.AssetRepo.GetAssetDetails(ctx, req.GetId())
// // 	if err != nil {
// // 		return nil, status.Errorf(codes.NotFound, "asset not found: %v", err)
// // 	}

// // 	// return ConvertToProtoAsset(asset), nil
// // 	return &pb_bdspro.AssetDetails{
// // 		Data: ConvertToProtoAsset(asset),
// // 	}, nil
// // }

// // func (BdsproGrpcServer) ShareAsset(context.Context, *pb_bdspro.ShareAssetRequest) (*pb_bdspro.Empty, error) {
// // 	return nil, status.Errorf(codes.Unimplemented, "method ShareAsset not implemented")
// // }
// // func (BdsproGrpcServer) SplitAsset(context.Context, *pb_bdspro.SplitAssetRequest) (*pb_bdspro.Empty, error) {
// // 	return nil, status.Errorf(codes.Unimplemented, "method SplitAsset not implemented")
// // }
// // func (BdsproGrpcServer) MergeAsset(context.Context, *pb_bdspro.MergeAssetRequest) (*pb_bdspro.Asset, error) {
// // 	return nil, status.Errorf(codes.Unimplemented, "method MergeAsset not implemented")
// // }
// // func (BdsproGrpcServer) GetAssetHistory(context.Context, *pb_bdspro.GetAssetRequest) (*pb_bdspro.AssetHistory, error) {
// // 	return nil, status.Errorf(codes.Unimplemented, "method GetAssetHistory not implemented")
// // }

// // func (BdsproGrpcServer) Ping(context.Context, *pb_bdspro.Empty) (*pb_bdspro.PingResponse, error) {
// // 	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
// // }
// // func (BdsproGrpcServer) ListProductInternal(context.Context, *pb_bdspro.ProductInternalRequest) (*pb_bdspro.ProductInternalResponse, error) {
// // 	return nil, status.Errorf(codes.Unimplemented, "method ListProductInternal not implemented")
// // }
// // func (BdsproGrpcServer) ConvertProductToAsset(context.Context, *pb_bdspro.CreateAssetRequestConvertProduct) (*pb_bdspro.Asset, error) {
// // 	return nil, status.Errorf(codes.Unimplemented, "method ConvertProductToAsset not implemented")
// // }
// // func (s BdsproGrpcServer) Archived(c context.Context, dto *pb_bdspro.ArchivedRequest) (*pb_bdspro.Empty, error) {
// // 	err := s.AssetUsecase.ArchiveAsset(c, dto.Id, dto.Archived)
// // 	return nil, err
// // }
