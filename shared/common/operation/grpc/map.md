info.FullMethod
→ "/types.tqd.MapWorkspaceService/CreateGeneratedReport"

operation.Code
→ "workspace.report.generate"

### OVERVIEW

gRPC request
│
▼
UnaryServerInterceptor
│
▼
Đọc info.FullMethod
│
▼
Tìm protobuf descriptor
│
▼
Đọc operation annotation
│
├── Không có ──────► handler(ctx, req)
│
▼
Có operation code
│
▼
operation.Bind(ctx, code)
│
├── Lỗi ───────────► trả gRPC Internal
│
▼
handler(bound, req)

### MORE DETAIL

grpc.UnaryServerInfo.FullMethod
"/types.tqd.MapWorkspaceService/CreateGeneratedReport"
│
▼
descriptorName()
"types.tqd.MapWorkspaceService.CreateGeneratedReport"
│
▼
protoregistry.GlobalFiles
│
▼
protoreflect.MethodDescriptor
│
▼
MethodOptions
│
▼
operationpb.E_Operation
│
▼
"workspace.report.generate"
│
▼
operation.Parse()
│
▼
operation.Bind(ctx, code)
│
▼
handler(boundCtx, req)
