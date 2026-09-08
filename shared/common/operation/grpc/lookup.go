package operationgrpc

import (
	"fmt"
	"strings"

	"common/operation"
	operationpb "pb/types/operation"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
)

// codeForMethod lấy operation code được khai báo trong protobuf của một gRPC method.
//
// fullMethod có dạng "/package.Service/Method".
// found = false khi method tồn tại nhưng không có operation annotation.
func codeForMethod(fullMethod string) (code operation.Code, found bool, err error) {
	name, err := descriptorName(fullMethod)
	if err != nil {
		return "", false, err
	}

	descriptor, err := protoregistry.GlobalFiles.FindDescriptorByName(name)
	if err != nil {
		return "", false, fmt.Errorf(
			"find protobuf descriptor for %q: %w",
			fullMethod,
			err,
		)
	}

	method, ok := descriptor.(protoreflect.MethodDescriptor)
	if !ok {
		return "", false, fmt.Errorf(
			"protobuf descriptor for %q is %T, want method",
			fullMethod,
			descriptor,
		)
	}

	options, ok := method.Options().(*descriptorpb.MethodOptions)
	if !ok || options == nil {
		return "", false, fmt.Errorf(
			"protobuf method options for %q are unavailable",
			fullMethod,
		)
	}

	return codeFromOptions(options)
}

// descriptorName chuyển gRPC full method name thành tên protobuf descriptor.
//
// Ví dụ:
//
//	"/package.Service/Method" → "package.Service.Method"
func descriptorName(fullMethod string) (protoreflect.FullName, error) {
	if !strings.HasPrefix(fullMethod, "/") {
		return "", fmt.Errorf(
			"gRPC full method %q must start with /",
			fullMethod,
		)
	}

	service, method, ok := strings.Cut(
		strings.TrimPrefix(fullMethod, "/"),
		"/",
	)
	if !ok ||
		service == "" ||
		method == "" ||
		strings.Contains(method, "/") {
		return "", fmt.Errorf(
			"gRPC full method %q must have /package.service/method form",
			fullMethod,
		)
	}

	return protoreflect.FullName(service + "." + method), nil
}

// codeFromOptions lấy operation annotation từ MethodOptions.
//
// Nếu RPC không có annotation, found = false.
// Nếu có annotation, code được parse và validate trước khi trả về.
func codeFromOptions(
	options *descriptorpb.MethodOptions,
) (code operation.Code, found bool, err error) {
	if options == nil ||
		!proto.HasExtension(options, operationpb.E_Operation) {
		return "", false, nil
	}

	value := proto.GetExtension(
		options,
		operationpb.E_Operation,
	)

	annotation, ok := value.(*operationpb.Operation)
	if !ok || annotation == nil {
		return "", false, fmt.Errorf(
			"operation annotation has type %T",
			value,
		)
	}

	code, err = operation.Parse(annotation.GetCode())
	if err != nil {
		return "", false, fmt.Errorf(
			"parse operation annotation %q: %w",
			annotation.GetCode(),
			err,
		)
	}

	return code, true, nil
}
