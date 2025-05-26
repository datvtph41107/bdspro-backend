package rest

import (
	_jwt "common/jwt"
	_routes "common/routes"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

type IRestService[M IRestModel, DTO any] interface {
	// RequiredRest(c *gin.Context, id uint64) (*RestModel[M], error)
	Create(c *gin.Context, md *RestModel[M]) (*RestModel[M], error)
	Update(c *gin.Context, id uint64, md *RestModel[M]) (*RestModel[M], error)
	Delete(c *gin.Context, id uint64) (bool, error)
	GetByID(c *gin.Context, id uint64) (*RestModel[M], error)
	Search(c *gin.Context, dto SearchDTO[DTO]) ([]RestModel[M], error)
	GenModel() RestModel[M]
}

type RestService[M IRestModel, DTO any] struct {
	RestRepo IRestRepo[M]
}

func (s RestService[M, DTO]) GenModel() RestModel[M] {
	var t RestModel[M]
	typ := reflect.TypeOf(t)

	// Tạo mới một con trỏ đến kiểu T
	val := reflect.New(typ).Interface()

	// *val
	// Ép kiểu con trỏ về *RestModel[M]
	r, ok := val.(RestModel[M])
	if !ok {
		fmt.Println("Failed to cast value to *RestModel[M]")
		return r
	}

	return r
}
func (s RestService[M, DTO]) Create(c *gin.Context, md *RestModel[M]) (*RestModel[M], error) {
	if err := s.RestRepo.Create(c, md); err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return md, nil
}

func (s *RestService[M, DTO]) Update(c *gin.Context, id uint64, md *RestModel[M]) (*RestModel[M], error) {
	if err := s.RestRepo.Update(c, id, md); err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return md, nil
}

func (s *RestService[M, DTO]) Delete(c *gin.Context, id uint64) (bool, error) {
	if err := s.RestRepo.Delete(c, id); err != nil {
		return false, &_routes.Except{
			Code:    500,
			Message: "Không thể xóa dữ liệu",
		}
	}

	return true, nil
}

func (s *RestService[M, DTO]) GetByID(c *gin.Context, id uint64) (*RestModel[M], error) {
	entity, err := s.RestRepo.GetByID(c, id)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (s *RestService[M, DTO]) Search(c *gin.Context, dto SearchDTO[DTO]) ([]RestModel[M], error) {
	profileId := _jwt.GetProfileId(c)
	skip := dto.Page * dto.Size
	limit := dto.Size
	if limit <= 0 {
		limit = 20
	}
	entities, err := s.RestRepo.Search(c,
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
