package product_private

import (
	product_accecss "bdspro/v1/mod/global/product_access"
	"bdspro/v1/models"
	"common/case/crud2"
	_db "common/db"
	"errors"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductPrivateRepo struct {
	crud2.BaseRepo[models.ProductPrivate]
	ProductAccess *product_accecss.ProductAccessRepo
}

func NewProductPrivateRepo(db *gorm.DB,
	ProductAccess *product_accecss.ProductAccessRepo,
) *ProductPrivateRepo {
	return &ProductPrivateRepo{
		BaseRepo:      crud2.BaseRepo[models.ProductPrivate]{DB: db},
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

func (repo *ProductPrivateRepo) GetPrivateFields(c *gin.Context, profileId uint64, product *models.Product) (*models.ProductPrivate, error) {
	privateData := repo.FindByProductId(&product.ID)
	if product.OwnerUserID == &profileId {
		return privateData, nil
	}
	access, _ := repo.ProductAccess.GetByUserAndProduct(profileId, product.ID)

	if access == nil {
		return nil, nil
	}
	// if err != nil {
	// 	return nil, err
	// }

	var result models.ProductPrivate
	FilterFields(privateData, access.GetFields(), &result)
	return &result, nil
}
func (r *ProductPrivateRepo) FindByProductId(productId *uint64) *models.ProductPrivate {
	var data *models.ProductPrivate

	_db.DB.
		Where("product_id = ? AND deleted_at IS NULL", productId).
		Order("created_at DESC").
		Limit(1).
		First(&data)

	return data
}

func (r *ProductPrivateRepo) UpdateInfo(c *gin.Context, productId uint64, entity *models.ProductPrivate) error {
	if entity == nil {
		return nil
	}
	data := r.FindByProductId(&productId)
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

	return _db.SaveWithAudit(c, &entity)
}
