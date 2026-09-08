package usecases

// type OrgMemAccessUsecase struct {
// 	repo        repo.SharingAccessRepo
// 	UProductUC  *UProductUsecase
// 	NotiService *providers.TransmitProvider
// }

// func NewOrgMemAccessUsecase(repo repo.SharingAccessRepo,
// 	UProductUC *UProductUsecase,
// 	NotiService *providers.TransmitProvider,
// ) *OrgMemAccessUsecase {
// 	return &OrgMemAccessUsecase{
// 		repo:        repo,
// 		UProductUC:  UProductUC,
// 		NotiService: NotiService,
// 	}
// }

// func (u *OrgMemAccessUsecase) GetProductAccessByOrgID(ctx context.Context, orgID int64, dto dto.ProductAccessSearch) ([]*domain.ProductAccess, error) {
// 	// return u.repo.GetProductAccessByOrgID(ctx, orgID, dto)
// 	return nil, nil
// }

// func (s *OrgMemAccessUsecase) HasPermission(c context.Context, id uint64) (bool, error) {
// 	// profileId := _utils.GetProfileIdWithContext(c)
// 	// product, err := s.repo.GetByTargetAndProduct(profileId, enums.ETargetTypeOrganization, id)
// 	// if err != nil {
// 	// 	return false, err
// 	// }

// 	// return product != nil, nil
// 	return true, nil
// }

// func (s *OrgMemAccessUsecase) BulkSave(c context.Context, productId uint64, dto dto.ProductAccessBulkRequest) (*[]domain.ProductAccess, error) {
// 	var entities []domain.ProductAccess
// 	_, err := s.HasPermission(c, productId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// todo: gọi sang org check thêm role nữa

// 	// đổi trạng thái sang nội bộ nếu riêng tư
// 	// if product.SaleVisibility == enums.EVisiblePrivate {
// 	// 	product.SaleVisibility = enums.EVisibleInternal
// 	// 	_db.SaveWithContext(c, product)
// 	// }
// 	if len(dto.Datas) > 0 {
// 		copier.Copy(&entities, &dto.Datas)
// 		for i := 0; i < len(entities); i++ {
// 			entities[i].ProductID = productId
// 			entities[i].TargetType = enums.EOwnerOfOrgnization
// 			entities[i].DeletedAt = nil
// 		}

// 		// assignments := append(
// 		// 	clause.AssignmentColumns([]string{"fields", "commission"}),
// 		// 	clause.Assignment{
// 		// 		Column: clause.Column{Name: "deleted_at"},
// 		// 		Value:  gorm.Expr("NULL"),
// 		// 	},
// 		// )

// 		// err := _db.DB.WithContext(c).Debug().
// 		// 	Clauses(clause.OnConflict{
// 		// 		Columns: []clause.Column{
// 		// 			{Name: "product_id"},
// 		// 			{Name: "target_id"},
// 		// 			{Name: "target_type"},
// 		// 		},
// 		// 		DoUpdates: clause.Set(assignments),
// 		// 	}).
// 		// 	Save(&entities).Error
// 		// if err != nil {
// 		// 	return nil, err
// 		// }
// 	}

// 	if len(dto.Deletes) > 0 {
// 		query := _db.DB.Model(&domain.ProductAccess{})
// 		for i, cond := range dto.Deletes {
// 			if i == 0 {
// 				query = query.Where("target_id = ? AND target_type = ?", cond.TargetID, cond.TargetType)
// 			} else {
// 				query = query.Or("target_id = ? AND target_type = ?", cond.TargetID, cond.TargetType)
// 			}
// 		}
// 		err := query.Update("deleted_at", time.Now()).Error
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	data := []*notificationpb.NotiNewRequest{}
// 	for i := 0; i < len(entities); i++ {
// 		data = append(data, &notificationpb.NotiNewRequest{
// 			Title:   "Chia sẻ sản phẩm",
// 			Message: "Bạn nhận được chia sẻ sản phẩm ",
// 			Type:    1,
// 			UserID:  entities[i].TargetId,
// 		})
// 		// entities[i].ProductID = id
// 	}

// 	s.NotiService.NotiClient.SendBatch(c, &notificationpb.NotiBatchRequest{
// 		Datas: data,
// 	})

// 	return &entities, nil
// }
