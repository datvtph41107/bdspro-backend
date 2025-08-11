package _owner

import (
	_jwt "common/jwt"
	_routes "common/routes"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

type IOwnerService[M IOwnerModel, DTO any] interface {
	RequiredOwner(c *gin.Context, id uint64) (*OwnerModel[M], error)
	Create(c *gin.Context, md *OwnerModel[M]) (*OwnerModel[M], error)
	Update(c *gin.Context, id uint64, md *OwnerModel[M]) (*OwnerModel[M], error)
	Delete(c *gin.Context, id uint64) (bool, error)
	GetByID(c *gin.Context, id uint64) (*OwnerModel[M], error)
	Search(c *gin.Context, dto SearchDTO[DTO]) ([]OwnerModel[M], error)
	GenModel() OwnerModel[M]
}

type OwnerService[M IOwnerModel, DTO any] struct {
	OwnerRepo IOwnerRepo[M]
}

func (s OwnerService[M, DTO]) GenModel() OwnerModel[M] {
	var t OwnerModel[M]
	typ := reflect.TypeOf(t)

	// Tạo mới một con trỏ đến kiểu T
	val := reflect.New(typ).Interface()

	// *val
	// Ép kiểu con trỏ về *OwnerModel[M]
	r, ok := val.(OwnerModel[M])
	if !ok {
		fmt.Println("Failed to cast value to *OwnerModel[M]")
		return r
	}

	return r
}
func (s OwnerService[M, DTO]) Create(c *gin.Context, md *OwnerModel[M]) (*OwnerModel[M], error) {
	profileId := _jwt.GetProfileId(c)
	md.OwnerID = profileId
	if err := s.OwnerRepo.Create(c, md); err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return md, nil
}

func (s *OwnerService[M, DTO]) RequiredOwner(c *gin.Context, id uint64) (*OwnerModel[M], error) {
	md, err := s.OwnerRepo.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	if md == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "ID không đúng",
		}
	}

	profileId := _jwt.GetProfileId(c)
	if md.OwnerID != profileId {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Bạn không có quyền truy cập",
		}
	}

	return md, nil
}
func (s *OwnerService[M, DTO]) Update(c *gin.Context, id uint64, md *OwnerModel[M]) (*OwnerModel[M], error) {
	if _, err := s.RequiredOwner(c, id); err != nil {
		return nil, err
	}
	if err := s.OwnerRepo.Update(c, id, md); err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return md, nil
}

func (s *OwnerService[M, DTO]) Delete(c *gin.Context, id uint64) (bool, error) {
	// id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	// if err != nil {
	// 	return false, &_routes.Except{
	// 		Code:    400,
	// 		Message: "ID không hợp lệ",
	// 	}
	// }
	if _, err := s.RequiredOwner(c, id); err != nil {
		return false, err
	}
	if err := s.OwnerRepo.Delete(c, id); err != nil {
		return false, &_routes.Except{
			Code:    500,
			Message: "Không thể xóa dữ liệu",
		}
	}

	return true, nil
}

func (s *OwnerService[M, DTO]) GetByID(c *gin.Context, id uint64) (*OwnerModel[M], error) {
	entity, err := s.RequiredOwner(c, id)
	if err != nil {
		return nil, err
	}
	// entity, err := s.OwnerRepo.GetByID(c, id)
	// if err != nil {
	// 	return nil, err
	// }

	return entity, err
}

func (s *OwnerService[M, DTO]) Search(c *gin.Context, dto SearchDTO[DTO]) ([]OwnerModel[M], error) {
	profileId := _jwt.GetProfileId(c)
	skip := dto.Page * dto.Size
	limit := dto.Size
	if limit <= 0 {
		limit = 20
	}
	entities, err := s.OwnerRepo.Search(c,
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
