package providers

import _usecase "common/domain/usecase"

func NewCodeDataUsecase() *_usecase.CodeDataUsecase {
	return _usecase.NewCodeUsecase()
}
