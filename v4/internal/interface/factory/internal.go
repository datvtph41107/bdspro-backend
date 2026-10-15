package factory

import authdto "userv4/modules/auth/dto"

type IFactory interface {
	GetProperties() authdto.PropertiesDTO
}
