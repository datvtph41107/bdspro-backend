package logic_crud

import (
	"context"
	"strconv"

	_routes "common/routes"

	"github.com/gin-gonic/gin"
)

type BaseUsecase[T any, R IRepo[T]] struct {
	Repo R
	Url  string
}

func (s *BaseUsecase[T, R]) Create(c *gin.Context) (*T, error) {
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

func (s *BaseUsecase[T, R]) Update(c *gin.Context) (*T, error) {
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

func (s *BaseUsecase[T, R]) Delete(c context.Context, id uint64) (bool, error) {
	// id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		return false, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		}
	}

	if err := s.Repo.Delete(c, id); err != nil {
		return false, &_routes.Except{
			Code:    500,
			Message: "Không thể xóa dữ liệu",
		}
	}

	return true, nil
}

func (s *BaseUsecase[T, R]) GetByID(c context.Context, id uint64) (*T, error) {
	// id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		}
	}

	entity, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, &_routes.Except{
			Code:    404,
			Message: "Không tìm thấy dữ liệu",
		}
	}

	return entity, nil
}

func (s *BaseUsecase[T, R]) GetAll(c context.Context) ([]T, error) {
	entities, err := s.Repo.GetAll()
	if err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return entities, nil
}
