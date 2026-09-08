package crud2

import (
	"strconv"

	_routes "common/routes"

	"github.com/gin-gonic/gin"
)

type BaseService[T any, R IBaseRepo[T]] struct {
	Repo R
}

func (s *BaseService[T, R]) Create(c *gin.Context) (*T, error) {
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Dữ liệu không hợp lệ",
		}
	}

	if err := s.Repo.Create(c, &entity); err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return &entity, nil
}

func (s *BaseService[T, R]) Update(c *gin.Context) (*T, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		}
	}

	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		return nil, &_routes.Except{
			Code: 400,
			// Message: "Dữ liệu không hợp lệ",
			Message: err.Error(),
		}
	}

	if err := s.Repo.Update(c, id, &entity); err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return &entity, nil
}

func (s *BaseService[T, R]) Delete(c *gin.Context) (bool, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return false, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		}
	}

	if err := s.Repo.Delete(c, id); err != nil {
		return false, &_routes.Except{
			Code:    500,
			Message: err.Error(), // "Không thể xóa dữ liệu",
		}
	}

	return true, nil
}

func (s *BaseService[T, R]) GetByID(c *gin.Context) (*T, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		}
	}

	entity, err := s.Repo.GetByID(c, id)
	if err != nil {
		return nil, &_routes.Except{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}

	return entity, nil
}

func (s *BaseService[T, R]) GetAll(c *gin.Context) ([]T, error) {
	entities, err := s.Repo.GetAll()
	if err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return entities, nil
}
