from pathlib import Path
import sys

if len(sys.argv) != 2:
    raise SystemExit("usage: materialize_s6_product_asset_unlink.py <source-root>")

root = Path(sys.argv[1]).resolve()


def replace_once(relative_path: str, old: str, new: str) -> None:
    path = root / relative_path
    text = path.read_text()
    count = text.count(old)
    if count != 1:
        raise SystemExit(f"{relative_path}: expected exactly one replacement target, got {count}")
    path.write_text(text.replace(old, new, 1))


replace_once(
    "bdspro-service/internal/domain/errors.go",
    '\tErrConcurrentUpdate = errors.New("ErrConcurrentUpdate")\n',
    '\tErrConcurrentUpdate            = errors.New("ErrConcurrentUpdate")\n'
    '\tErrProductAssetLinkCheckFailed = errors.New("product asset link check failed")\n'
    '\tErrProductAssetLinkNotFound    = errors.New("product asset link not found")\n'
    '\tErrProductAssetUnlinkFailed    = errors.New("product asset unlink failed")\n',
)

replace_once(
    "bdspro-service/internal/usecases/product_asset_usecase.go",
    'import (\n\t"bdspro/internal/dto"\n',
    'import (\n\t"bdspro/internal/domain"\n\t"bdspro/internal/dto"\n',
)

old_unlink = '''func (uc *ProductAssetUsecase) UnlinkProductAsset(ctx context.Context, req *dto.UnlinkProductAssetRequest) error {
\t// Kiểm tra liên kết tồn tại
\texists, err := uc.ProductAssetRepo.CheckExists(ctx, req.ProductID, req.AssetID)
\tif err != nil {
\t\treturn &_routes.Except{
\t\t\tCode:    500,
\t\t\tMessage: "Lỗi kiểm tra liên kết",
\t\t}
\t}
\tif !exists {
\t\treturn &_routes.Except{
\t\t\tCode:    404,
\t\t\tMessage: "Liên kết không tồn tại",
\t\t}
\t}

\t// Hủy liên kết
\terr = uc.ProductAssetRepo.Unlink(ctx, req.ProductID, req.AssetID)
\tif err != nil {
\t\treturn &_routes.Except{
\t\t\tCode:    500,
\t\t\tMessage: "Lỗi hủy liên kết",
\t\t}
\t}

\t// Cập nhật AssetCount cho product sau khi hủy liên kết
\tuc.updateAssetCountForProducts(ctx, []uint64{req.ProductID})

\treturn nil
}
'''

new_unlink = '''func (uc *ProductAssetUsecase) UnlinkProductAsset(ctx context.Context, req *dto.UnlinkProductAssetRequest) error {
\t// Kiểm tra liên kết tồn tại
\texists, err := uc.ProductAssetRepo.CheckExists(ctx, req.ProductID, req.AssetID)
\tif err != nil {
\t\treturn fmt.Errorf("%w: %v", domain.ErrProductAssetLinkCheckFailed, err)
\t}
\tif !exists {
\t\treturn domain.ErrProductAssetLinkNotFound
\t}

\t// Hủy liên kết
\terr = uc.ProductAssetRepo.Unlink(ctx, req.ProductID, req.AssetID)
\tif err != nil {
\t\treturn fmt.Errorf("%w: %v", domain.ErrProductAssetUnlinkFailed, err)
\t}

\t// Cập nhật AssetCount cho product sau khi hủy liên kết
\tuc.updateAssetCountForProducts(ctx, []uint64{req.ProductID})

\treturn nil
}
'''

replace_once(
    "bdspro-service/internal/usecases/product_asset_usecase.go",
    old_unlink,
    new_unlink,
)

replace_once(
    "bdspro-service/infra/handler/product_handler.go",
    '\t"bdspro/internal/dto"\n',
    '\t"bdspro/internal/domain"\n\t"bdspro/internal/dto"\n',
)
replace_once(
    "bdspro-service/infra/handler/product_handler.go",
    '\t"context"\n',
    '\t"context"\n\t"errors"\n',
)
replace_once(
    "bdspro-service/infra/handler/product_handler.go",
    '''\terr := s.ProductAssetUC.UnlinkProductAsset(ctx, unlinkReq)
\tif err != nil {
\t\treturn nil, err
\t}
''',
    '''\terr := s.ProductAssetUC.UnlinkProductAsset(ctx, unlinkReq)
\tif err != nil {
\t\treturn nil, mapProductAssetUnlinkError(err)
\t}
''',
)

marker = '\n// @Summary Lấy danh sách tài sản theo sản phẩm\n'
mapper = '''
func mapProductAssetUnlinkError(err error) error {
\tif err == nil {
\t\treturn nil
\t}

\tswitch {
\tcase errors.Is(err, domain.ErrProductAssetLinkCheckFailed):
\t\treturn status.Error(codes.Unknown, "500: Lỗi kiểm tra liên kết")
\tcase errors.Is(err, domain.ErrProductAssetLinkNotFound):
\t\treturn status.Error(codes.Unknown, "404: Liên kết không tồn tại")
\tcase errors.Is(err, domain.ErrProductAssetUnlinkFailed):
\t\treturn status.Error(codes.Unknown, "500: Lỗi hủy liên kết")
\tdefault:
\t\treturn err
\t}
}
'''
replace_once(
    "bdspro-service/infra/handler/product_handler.go",
    marker,
    '\n' + mapper + marker,
)

(root / "bdspro-service/infra/handler/product_asset_unlink_error_test.go").write_text(
    '''package handler

import (
\t"errors"
\t"fmt"
\t"testing"

\t"bdspro/internal/domain"
\t"google.golang.org/grpc/codes"
\t"google.golang.org/grpc/status"
)

func TestMapProductAssetUnlinkErrorLegacyContract(t *testing.T) {
\tcases := []struct {
\t\tname        string
\t\terr         error
\t\twantMessage string
\t}{
\t\t{"check failure", fmt.Errorf("repo check: %w", domain.ErrProductAssetLinkCheckFailed), "500: Lỗi kiểm tra liên kết"},
\t\t{"not found", domain.ErrProductAssetLinkNotFound, "404: Liên kết không tồn tại"},
\t\t{"unlink failure", fmt.Errorf("repo unlink: %w", domain.ErrProductAssetUnlinkFailed), "500: Lỗi hủy liên kết"},
\t}

\tfor _, tc := range cases {
\t\tt.Run(tc.name, func(t *testing.T) {
\t\t\tgot := status.Convert(mapProductAssetUnlinkError(tc.err))
\t\t\tif got.Code() != codes.Unknown {
\t\t\t\tt.Fatalf("code=%s want=%s", got.Code(), codes.Unknown)
\t\t\t}
\t\t\tif got.Message() != tc.wantMessage {
\t\t\t\tt.Fatalf("message=%q want=%q", got.Message(), tc.wantMessage)
\t\t\t}
\t\t\tif len(got.Details()) != 0 {
\t\t\t\tt.Fatalf("details=%d want=0", len(got.Details()))
\t\t\t}
\t\t})
\t}
}

func TestMapProductAssetUnlinkErrorPassthrough(t *testing.T) {
\tother := errors.New("unrelated")
\tif got := mapProductAssetUnlinkError(other); got != other {
\t\tt.Fatalf("got=%v want same unrelated error", got)
\t}
\tif got := mapProductAssetUnlinkError(nil); got != nil {
\t\tt.Fatalf("got=%v want nil", got)
\t}
}
'''
)
