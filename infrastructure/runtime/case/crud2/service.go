package crud2

import (
	"fmt"
	"strconv"

	_errors "common/errors"

	"github.com/gin-gonic/gin"
)

type BaseService[T any, R IBaseRepo[T]] struct {
	Repo R
}

func (s *BaseService[T, R]) Create(c *gin.Context) (*T, error) {
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		return nil, _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithCause(err))
	}
	if err := s.Repo.Create(c, &entity); err != nil {
		return nil, fmt.Errorf("create record: %w", err)
	}
	return &entity, nil
}

func (s *BaseService[T, R]) Update(c *gin.Context) (*T, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, _errors.ReturnError(_errors.ResourceIDInvalid, _errors.WithCause(err))
	}

	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		return nil, _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithCause(err))
	}
	if err := s.Repo.Update(c, id, &entity); err != nil {
		return nil, fmt.Errorf("update record %d: %w", id, err)
	}
	return &entity, nil
}

func (s *BaseService[T, R]) Delete(c *gin.Context) (bool, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return false, _errors.ReturnError(_errors.ResourceIDInvalid, _errors.WithCause(err))
	}
	if err := s.Repo.Delete(c, id); err != nil {
		return false, fmt.Errorf("delete record %d: %w", id, err)
	}
	return true, nil
}

func (s *BaseService[T, R]) GetByID(c *gin.Context) (*T, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, _errors.ReturnError(_errors.ResourceIDInvalid, _errors.WithCause(err))
	}

	entity, err := s.Repo.GetByID(c, id)
	if err != nil {
		return nil, _errors.ReturnError(_errors.DataNotFound, _errors.WithCause(err))
	}
	return entity, nil
}

func (s *BaseService[T, R]) GetAll(c *gin.Context) ([]T, error) {
	entities, err := s.Repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("list records: %w", err)
	}
	return entities, nil
}
