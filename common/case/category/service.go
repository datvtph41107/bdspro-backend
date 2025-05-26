package _category

import (
	_jwt "common/jwt"
	_routes "common/routes"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

type ICategoryService[M ICategoryModel, DTO any] interface {
	RequiredCategory(c *gin.Context, id uint64) (*CategoryModel[M], error)
	Create(c *gin.Context, md *CategoryModel[M]) (*CategoryModel[M], error)
	Update(c *gin.Context, id uint64, md *CategoryModel[M]) (*CategoryModel[M], error)
	Delete(c *gin.Context, id uint64) (bool, error)
	GetByID(c *gin.Context, id uint64) (*CategoryModel[M], error)
	Search(c *gin.Context, dto SearchDTO[DTO]) ([]CategoryModel[M], error)
	GenModel() CategoryModel[M]
}

type CategoryService[M ICategoryModel, DTO any] struct {
	CategoryRepo ICategoryRepo[M]
}

func (s CategoryService[M, DTO]) GenModel() CategoryModel[M] {
	var t CategoryModel[M]
	typ := reflect.TypeOf(t)

	// Tạo mới một con trỏ đến kiểu T
	val := reflect.New(typ).Interface()

	// *val
	// Ép kiểu con trỏ về *CategoryModel[M]
	r, ok := val.(CategoryModel[M])
	if !ok {
		fmt.Println("Failed to cast value to *CategoryModel[M]")
		return r
	}

	return r
}
func (s CategoryService[M, DTO]) Create(c *gin.Context, md *CategoryModel[M]) (*CategoryModel[M], error) {
	md.IsCustom = false
	if err := s.CategoryRepo.Create(c, md); err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return md, nil
}

func (s *CategoryService[M, DTO]) RequiredCategory(c *gin.Context, id uint64) (*CategoryModel[M], error) {
	md, err := s.CategoryRepo.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	if md == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "ID không đúng",
		}
	}

	// profileId := _jwt.GetProfileId(c)
	// if md.CategoryID != profileId {
	// 	return nil, &_routes.Except{
	// 		Code:    400,
	// 		Message: "Bạn không có quyền truy cập",
	// 	}
	// }

	return md, nil
}
func (s *CategoryService[M, DTO]) Update(c *gin.Context, id uint64, md *CategoryModel[M]) (*CategoryModel[M], error) {
	if _, err := s.RequiredCategory(c, id); err != nil {
		return nil, err
	}
	if err := s.CategoryRepo.Update(c, id, md); err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return md, nil
}

func (s *CategoryService[M, DTO]) Delete(c *gin.Context, id uint64) (bool, error) {
	// id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	// if err != nil {
	// 	return false, &_routes.Except{
	// 		Code:    400,
	// 		Message: "ID không hợp lệ",
	// 	}
	// }
	if _, err := s.RequiredCategory(c, id); err != nil {
		return false, err
	}
	if err := s.CategoryRepo.Delete(c, id); err != nil {
		return false, &_routes.Except{
			Code:    500,
			Message: "Không thể xóa dữ liệu",
		}
	}

	return true, nil
}

func (s *CategoryService[M, DTO]) GetByID(c *gin.Context, id uint64) (*CategoryModel[M], error) {
	entity, err := s.RequiredCategory(c, id)
	if err != nil {
		return nil, err
	}
	// entity, err := s.CategoryRepo.GetByID(c, id)
	// if err != nil {
	// 	return nil, err
	// }

	return entity, err
}

func (s *CategoryService[M, DTO]) Search(c *gin.Context, dto SearchDTO[DTO]) ([]CategoryModel[M], error) {
	profileId := _jwt.GetProfileId(c)
	skip := dto.Page * dto.Size
	limit := dto.Size
	if limit <= 0 {
		limit = 20
	}
	entities, err := s.CategoryRepo.Search(c,
		profileId,
		skip,
		limit,
	)
	if err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return entities, nil
}
