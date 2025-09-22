package postgres

import (
	"bdspro/internal/domain"
	"common/case/crud3"
	"context"
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.ProductPrivateRepo
type GormProductPrivateRepo struct {
	crud3.BaseRepo[domain.ProductPrivate]
	ProductAccess *SharingAccessPostgre
}

func NewGormProductPrivateRepo(db *gorm.DB,
	ProductAccess *SharingAccessPostgre,
) *GormProductPrivateRepo {
	return &GormProductPrivateRepo{
		BaseRepo:      crud3.BaseRepo[domain.ProductPrivate]{DB: db},
		ProductAccess: ProductAccess,
	}
}

func FilterFields(obj interface{}, fields []string, desc interface{}) error {
	// Lấy reflect.Value của obj và desc
	objVal := reflect.ValueOf(obj)
	descVal := reflect.ValueOf(desc)

	// Kiểm tra nếu desc không phải là con trỏ hoặc không trỏ đến struct
	if descVal.Kind() != reflect.Ptr || descVal.Elem().Kind() != reflect.Struct {
		return errors.New("desc must be a pointer to a struct")
	}

	// Lấy giá trị thực tế từ con trỏ
	descVal = descVal.Elem()
	objType := reflect.TypeOf(obj)

	// Nếu obj là con trỏ, lấy giá trị thực tế của nó
	if objVal.Kind() == reflect.Ptr {
		objVal = objVal.Elem()
		objType = objType.Elem()
	}

	// Duyệt qua danh sách fields để sao chép giá trị
	for _, field := range fields {
		// Lấy thông tin của field trong obj
		objField, found := objType.FieldByName(field)
		if !found {
			continue
		}

		// Lấy giá trị của field trong obj
		fieldValue := objVal.FieldByName(field)

		// Kiểm tra nếu desc cũng có field này
		descField := descVal.FieldByName(field)
		if descField.IsValid() && descField.CanSet() {
			// Gán giá trị từ obj sang desc
			descField.Set(fieldValue)
		} else {
			fmt.Printf("Field %s không thể set vào struct đích\n", objField.Name)
		}
	}

	return nil
}

func (repo *GormProductPrivateRepo) GetPrivateFields(c context.Context, profileId uint64, product *domain.Product) (*domain.ProductPrivate, error) {
	privateData := repo.FindByProductId(c, &product.ID)
	if product.OwnerID == profileId {
		return privateData, nil
	}
	access, _ := repo.ProductAccess.GetByUserAndProduct(profileId, product.ID)

	if access == nil {
		return nil, nil
	}
	// if err != nil {
	// 	return nil, err
	// }

	var result domain.ProductPrivate
	FilterFields(privateData, access.GetFields(), &result)
	return &result, nil
}
func (r *GormProductPrivateRepo) FindByProductId(c context.Context, productId *uint64) *domain.ProductPrivate {
	var data *domain.ProductPrivate

	GetDB(c, r.DB).
		Where("product_id = ? AND deleted_at IS NULL", productId).
		Order("created_at DESC").
		Limit(1).
		First(&data)

	return data
}

func (r *GormProductPrivateRepo) UpdateInfo(c context.Context, productId uint64, entity *domain.ProductPrivate) error {
	if entity == nil {
		return nil
	}
	data := r.FindByProductId(c, &productId)
	// if oldPrice != nil {
	// 	diff := !_utils.CompareEqual(&oldPrice.Currency, &entity.Currency) ||
	// 		!_utils.CompareEqual(&oldPrice.PriceOwner, &entity.PriceOwner) ||
	// 		!_utils.CompareEqual(&oldPrice.SalePrice, &entity.SalePrice) ||
	// 		!_utils.CompareEqual(&oldPrice.SaleCommission, &entity.SaleCommission) ||
	// 		!_utils.CompareEqual(&oldPrice.RentPrice, &entity.RentPrice) ||
	// 		!_utils.CompareEqual(&oldPrice.RentCommission, &entity.RentCommission) ||
	// 		!_utils.CompareEqual(&oldPrice.RentPaymentCycle, &entity.RentPaymentCycle)
	// 		// oldPrice.SaleCommissionType != entity.Price.SaleCommissionType ||
	// 		// oldPrice.RentCommissionType != entity.Price.RentCommissionType ||
	// 	if !diff {
	// 		return nil
	// 	}
	// }
	if data != nil {
		entity.ID = data.ID
		// return _db.SaveWithAudit(c, &entity)
	}
	entity.ProductId = productId

	return GetDB(c, r.DB).Save(&entity).Error
}
