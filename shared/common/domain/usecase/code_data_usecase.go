package _usecase

import (
	"context"
	"fmt"
)

type CodeDataUsecase struct {
}

func NewCodeUsecase() *CodeDataUsecase {
	return &CodeDataUsecase{}
}
func (uc *CodeDataUsecase) GetProductCode(c context.Context, productId uint64) string {
	return fmt.Sprintf("SP%06d", productId)
}
func (uc *CodeDataUsecase) GetProductCodePtr(c context.Context, productId *uint64) *string {
	if productId == nil {
		return nil
	}
	code := fmt.Sprintf("SP%06d", *productId)
	return &code
}
func (uc *CodeDataUsecase) GetAssetCode(c context.Context, assetId uint64) string {
	return fmt.Sprintf("AS%06d", assetId)
}
func (uc *CodeDataUsecase) GetAssetCodePtr(c context.Context, assetId *uint64) *string {
	if assetId == nil {
		return nil
	}
	code := fmt.Sprintf("AS%06d", *assetId)
	return &code
}
func (uc *CodeDataUsecase) GetPostCodePtr(c context.Context, postId *uint64) *string {
	if postId == nil {
		return nil
	}
	code := fmt.Sprintf("PO%06d", *postId)
	return &code
}
func (uc *CodeDataUsecase) GetPostCode(c context.Context, postId uint64) string {
	return fmt.Sprintf("PO%06d", postId)
}
func (uc *CodeDataUsecase) GetPostCodes(c context.Context, postIds []uint64) []string {
	postCodes := make([]string, len(postIds))
	for i, postId := range postIds {
		postCodes[i] = fmt.Sprintf("PO%06d", postId)
	}
	return postCodes
}
