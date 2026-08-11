# Current -> target

| Current source | Target role | Action khi migrate |
|---|---|---|
| `requestcontext` | request identity | keep/refine |
| giant `grpc_metadata_interceptor.go` | nhiều việc trộn | target RPC chuyển sang `requestgrpc` + `identitygrpc`, sau đó retire dần |
| `quotasubject*` ở ingress | xác định payer quá sớm | không dùng trong target report; access service quyết subject |
| `operationbinding*`, `operationruntime` | operation cũ quá nhiều field | target report dùng `businessoperation` + `grpcoperation` nhỏ |
| legacy quota `userID/service/action` | tự suy identity/action | target report skip hoàn toàn |
| `RedisQuotaStorage` cũ | quota runtime cũ | target report dùng `infra/redisquota` |
| `QHUserReported.ClientRequestID` | chống report trùng | giữ làm command key bridge, widen 128 |
| catalog V2 | source cho access | dùng ở G4 target sau DB migration check |
| TQD -> file-service HTTP client | service-to-service HTTP | thay bằng `infra/filegrpc` |
| file-service Gin business methods | direct HTTP service | giữ legacy trong migration; target internal API là gRPC |
