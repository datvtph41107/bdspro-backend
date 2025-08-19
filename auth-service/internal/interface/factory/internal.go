package factory

import "auth/internal/dto"

type IFactory interface {
	GetProperties() dto.PropertiesDTO
}
