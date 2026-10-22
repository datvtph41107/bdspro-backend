package quota_test

import (
	"mdn4/modules/tqd/services/quota"
	"testing"

	operationpb "pb/types/operation"
	tqdpb "pb/types/tqd"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// So với danh sách tự chép trong test, annotation protobuf là bằng chứng độc
// lập: thêm RPC tính quota mà quên usage projection phải làm test thất bại.
func TestProfileOperationsMatchRPCContract(t *testing.T) {
	annotated := map[string]bool{}
	for _, file := range []protoreflect.FileDescriptor{tqdpb.File_tqd_parcel_proto, tqdpb.File_tqd_map_workspace_proto} {
		services := file.Services()
		for i := 0; i < services.Len(); i++ {
			methods := services.Get(i).Methods()
			for j := 0; j < methods.Len(); j++ {
				options := methods.Get(j).Options()
				if !proto.HasExtension(options, operationpb.E_Operation) {
					continue
				}
				annotation := proto.GetExtension(options, operationpb.E_Operation).(*operationpb.Operation)
				annotated[annotation.GetCode()] = true
			}
		}
	}
	seen := map[string]bool{}
	for _, code := range quota.ProfileOperations() {
		if !code.IsValid() || seen[string(code)] || !annotated[string(code)] {
			t.Fatalf("operation sai, trùng hoặc thiếu trong protobuf: %q", code)
		}
		seen[string(code)] = true
	}
	for code := range annotated {
		if !seen[code] {
			t.Errorf("RPC operation %q chưa có trong usage projection", code)
		}
	}
}

func TestProfileOperationsCannotBeMutatedByCaller(t *testing.T) {
	first := quota.ProfileOperations()
	if len(first) == 0 {
		t.Fatal("usage registry is empty")
	}
	original := first[0]
	first[0] = "caller.modified.operation"
	if quota.ProfileOperations()[0] != original {
		t.Fatal("caller đã làm thay đổi registry dùng chung")
	}
}
