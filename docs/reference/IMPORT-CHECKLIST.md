# Import Checklist

## 1. Extract

Giải nén candidate vào workspace local. Archive không chứa `.env`, `bin`,
`.tmp`, logs hoặc runtime cache.

## 2. Prepare

```bash
make setup
make help
```

`setup` tạo root/service `.env` local và bootstrap pinned toolchain.

## 3. Source proof

```bash
make verify
```

Nếu chỉ thay migration:

```bash
make verify-migrations
```

## 4. Daily work on one service

```bash
make deps service=payment-service
make dev service=payment-service
make test service=payment-service
```

`make dev` là canonical native/Air loop. Logs của process đang sửa nằm ngay terminal/VS Code task.

## 5. Whole-system integration

```bash
make up
make status
make smoke
```

Chỉ rebuild image khi cần proof packaging/integration source mới:

```bash
make rebuild service=payment-service
make rebuild   # full candidate
```

Container logs của integration stack:

```bash
make logs service=payment-service
```

## 6. Migration

```bash
make migrate service=payment-service
make migrate
```

Không dùng `make reset` nếu cần giữ local data. `make reset` là destructive và
có confirmation phrase.

## 7. Final local acceptance

```bash
make accept
```

Nếu fail, giữ nguyên output/log và sửa theo owner service; không sửa production
deployment files để chữa lỗi local developer workflow.
