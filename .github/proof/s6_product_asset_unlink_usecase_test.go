package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	"context"
	"errors"
	"testing"
)

type productAssetUnlinkProofRepo struct {
	repo.ProductAssetRepo
	exists    bool
	checkErr  error
	unlinkErr error
}

func (r *productAssetUnlinkProofRepo) CheckExists(context.Context, uint64, uint64) (bool, error) {
	return r.exists, r.checkErr
}

func (r *productAssetUnlinkProofRepo) Unlink(context.Context, uint64, uint64) error {
	return r.unlinkErr
}

func TestProofProductAssetUnlinkSemanticErrors(t *testing.T) {
	req := &dto.UnlinkProductAssetRequest{ProductID: 11, AssetID: 22}
	checkCause := errors.New("check failed")
	unlinkCause := errors.New("unlink failed")

	cases := []struct {
		name string
		repo *productAssetUnlinkProofRepo
		want error
	}{
		{
			name: "check failure",
			repo: &productAssetUnlinkProofRepo{checkErr: checkCause},
			want: domain.ErrProductAssetLinkCheckFailed,
		},
		{
			name: "link not found",
			repo: &productAssetUnlinkProofRepo{exists: false},
			want: domain.ErrProductAssetLinkNotFound,
		},
		{
			name: "unlink failure",
			repo: &productAssetUnlinkProofRepo{exists: true, unlinkErr: unlinkCause},
			want: domain.ErrProductAssetUnlinkFailed,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			uc := &ProductAssetUsecase{ProductAssetRepo: tc.repo}
			err := uc.UnlinkProductAsset(context.Background(), req)
			if !errors.Is(err, tc.want) {
				t.Fatalf("err=%v want errors.Is(..., %v)", err, tc.want)
			}
		})
	}
}
