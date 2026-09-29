package factory

import "user/internal/dto"

type IFactory interface {
	GetProperties() dto.PropertiesDTO
}
