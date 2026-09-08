package postgres

// type AdminProductRepo struct {
// 	crud2.BaseRepo[domain.Product]
// 	PriceRepo      *product_price.ProductPriceRepo
// 	ProductAccess  *product_accecss.ProductAccessRepo
// 	ProductPrivate *product_private.ProductPrivateRepo
// }

// func NewAdminProductRepo(db *gorm.DB,
// 	PriceRepo *product_price.ProductPriceRepo,
// 	ProductAccess *product_accecss.ProductAccessRepo,
// 	ProductPrivate *product_private.ProductPrivateRepo,
// ) *AdminProductRepo {
// 	return &AdminProductRepo{
// 		BaseRepo:       crud2.BaseRepo[domain.Product]{DB: db},
// 		PriceRepo:      PriceRepo,
// 		ProductAccess:  ProductAccess,
// 		ProductPrivate: ProductPrivate,
// 	}
// }

// // func (r *AdminProductRepo) Create(c *gin.Context, entity *domain.Product) error {
// // }

// // Tạo mới bản ghi
// func (r *AdminProductRepo) Create(c *gin.Context, entity *domain.Product) error {
// 	err := r.DB.WithContext(c).Create(entity).Error
// 	if err == nil {
// 		err = r.PriceRepo.UpdatePrice(c, *entity.LastPriceID, entity.Price)
// 	}
// 	return err
// }

// // Cập nhật bản ghi
// func (r *AdminProductRepo) Update(c *gin.Context, id uint64, entity *domain.Product) error {
// 	price := entity.Price
// 	privateInfo := entity.PrivateData
// 	entity.Price = nil
// 	entity.PrivateData = nil

// 	err := r.PriceRepo.UpdatePrice(c, id, price)
// 	if err != nil {
// 		return err
// 	}

// 	entity.LastPriceID = &price.ID
// 	err = r.DB.WithContext(c).Model(entity).
// 		Where("id = ? and deleted_at is null", id).
// 		Updates(entity).Error

// 	if err != nil {
// 		return err
// 	}
// 	err = r.ProductPrivate.UpdateInfo(c, id, privateInfo)

// 	return err
// }

// func (r *AdminProductRepo) GetAll() ([]domain.Product, error) {
// 	var productsWithPrices []domain.Product

// 	r.DB.
// 		Preload("MediaList").
// 		Preload("Price", func(db *gorm.DB) *gorm.DB {
// 			return db.Order("product_price.created_at asc")
// 		}).
// 		Find(&productsWithPrices)

// 	return productsWithPrices, nil
// }

// // Lấy bản ghi theo ID (bỏ qua những bản ghi đã bị xóa mềm)
// func (r *AdminProductRepo) GetByID(c *gin.Context, id uint64) (*domain.Product, error) {
// 	log.Printf("LOAD_DETAIL", id)
// 	// Table("products"). // Bỏ alias 'p' vì GORM không hỗ trợ alias trong .Table()
// 	// Select("products.id, products.name, lp.price").
// 	// Select("*").
// 	var entity *domain.Product
// 	price := r.PriceRepo.FindByProductId(&id)

// 	err := _db.DB.
// 		Table("products p").
// 		Joins(`LEFT JOIN product_price lp ON lp.product_id = ?
//         	AND lp.created_at = (SELECT MAX(created_at) FROM product_price p2 WHERE p2.product_id = lp.product_id)`, id).
// 		Preload("PropertyType").
// 		Preload("DocType").
// 		Preload("Amenities", "deleted_at is null").
// 		Preload("Apartment", "deleted_at is null").
// 		Preload("MediaList", "deleted_at is null").
// 		Preload("HouseInfo", "deleted_at is null").
// 		Preload("CreatedUser", "deleted_at is null").
// 		Preload("UpdatedUser", "deleted_at is null").
// 		// Preload("PrivateData").
// 		Where("p.id = ? AND p.deleted_at IS NULL", id).
// 		First(&entity).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	entity.Price = price

// 	if entity.Apartment != nil {
// 		var build *domain.ProjectBuild
// 		_ = _db.DB.
// 			Where("id = ? AND deleted_at IS NULL", entity.Apartment.BuildID).
// 			First(&build).Error
// 		entity.Build = build

// 		var attribute *domain.ApartmentAttribute
// 		_ = _db.DB.
// 			Where("id = ? AND deleted_at IS NULL", entity.Apartment.ApartmentAttrID).
// 			First(&attribute).Error
// 		entity.ApartmentAttr = attribute
// 	}

// 	if entity.Build != nil {
// 		var project *domain.Project
// 		_ = _db.DB.
// 			Where("id = ? AND deleted_at IS NULL", entity.Build.ProjectID).
// 			First(&project).Error
// 		entity.Project = project
// 	}

// 	if entity.Project != nil {
// 		var developer *domain.Developer
// 		_ = _db.DB.
// 			Where("id = ? AND deleted_at IS NULL", entity.Project.DeveloperID).
// 			First(&developer).Error
// 		entity.Developer = developer
// 	}

// 	profileId := _jwt.GetProfileId(c)

// 	privateData, _ := r.ProductPrivate.GetPrivateFields(c, profileId, entity)
// 	if privateData != nil {
// 		entity.PrivateData = privateData
// 	}

// 	return entity, nil
// }
