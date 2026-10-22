package factory

import authdto "mdn4/modules/auth/dto"

type IFactory interface {
	GetProperties() authdto.PropertiesDTO
}
