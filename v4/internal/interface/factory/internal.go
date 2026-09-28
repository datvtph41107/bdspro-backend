package factory

import "userv4/internal/dto"

type IFactory interface {
	GetProperties() dto.PropertiesDTO
}
