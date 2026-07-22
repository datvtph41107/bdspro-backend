package flow

// import (
//     "context"
//     "database/sql"
//     "fmt"
//     "time"
//     "gorm.io/gorm"
//     apperr "bdspro/platform/errors"
// )

// // TxOption — functional options cho transaction
// type TxOption func(*txConfig)
// type txConfig struct {
//     isolation sql.IsolationLevel
//     timeout   time.Duration
// }

// func WithIsolation(level sql.IsolationLevel) TxOption {
//     return func(c *txConfig) { c.isolation = level }
// }
// func WithTimeout(d time.Duration) TxOption {
//     return func(c *txConfig) { c.timeout = d }
// }

// // TransactionProvider — interface, không phải concrete type
// // Usecase nhận interface này — không biết gorm tồn tại
// type TransactionProvider interface {
//     WithTransaction(ctx context.Context, fn func(ctx context.Context) error, opts ...TxOption) error
// }

// type gormTxProvider struct{ db *gorm.DB }

// func NewGormTxProvider(db *gorm.DB) TransactionProvider {
//     return &gormTxProvider{db: db}
// }

// func (p *gormTxProvider) WithTransaction(
//     ctx context.Context,
//     fn func(ctx context.Context) error,
//     opts ...TxOption,
// ) error {
//     cfg := &txConfig{
//         isolation: sql.LevelReadCommitted,
//         timeout:   10 * time.Second,
//     }
//     for _, o := range opts { o(cfg) }

//     ctx, cancel := context.WithTimeout(ctx, cfg.timeout)
//     defer cancel()

//     return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
//         // Inject tx vào ctx — repo dùng dbFromCtx() để lấy active tx
//         txCtx := context.WithValue(ctx, txKey{}, tx)
//         if err := fn(txCtx); err != nil {
//             return err // gorm auto rollback
//         }
//         return nil // gorm auto commit
//     }, &sql.TxOptions{Isolation: cfg.isolation})
// }

// type txKey struct{}

// // DBFromCtx — repo dùng để lấy active transaction nếu có
// // Nếu không có tx trong ctx → dùng fallback DB connection
// func DBFromCtx(ctx context.Context, fallback *gorm.DB) *gorm.DB {
//     if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
//         return tx.WithContext(ctx)
//     }
//     return fallback.WithContext(ctx)
// }
// // platform/core/flow/idempotency.go
// package flow

// import (
//     "context"
//     "encoding/json"
//     "errors"
//     "fmt"
//     "time"
//     "github.com/redis/go-redis/v9"
//     apperr "bdspro/platform/errors"
// )

// var ErrInProgress = errors.New("request in progress")

// // IdempotencyStore — interface, bất kỳ backend nào implement được
// type IdempotencyStore interface {
//     Acquire(ctx context.Context, key string, ttl time.Duration) (acquired bool, err error)
//     GetResult(ctx context.Context, key string, dest any) (found bool, err error)
//     SetResult(ctx context.Context, key string, result any, ttl time.Duration) error
//     Release(ctx context.Context, key string) error
// }

// type redisIdempotency struct{ client *redis.Client }

// func NewRedisIdempotency(client *redis.Client) IdempotencyStore {
//     return &redisIdempotency{client: client}
// }

// const inProgressMarker = "__in_progress__"

// func (s *redisIdempotency) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
//     // SET NX → atomic, không có race condition
//     return s.client.SetNX(ctx, key, inProgressMarker, ttl).Result()
// }

// func (s *redisIdempotency) GetResult(ctx context.Context, key string, dest any) (bool, error) {
//     data, err := s.client.Get(ctx, key).Bytes()
//     if errors.Is(err, redis.Nil) { return false, nil }
//     if err != nil { return false, fmt.Errorf("idempotency.get: %w", err) }
//     if string(data) == inProgressMarker {
//         return false, ErrInProgress // concurrent same request
//     }
//     return true, json.Unmarshal(data, dest)
// }

// func (s *redisIdempotency) SetResult(ctx context.Context, key string, result any, ttl time.Duration) error {
//     data, err := json.Marshal(result)
//     if err != nil { return err }
//     return s.client.Set(ctx, key, data, ttl).Err()
// }

// func (s *redisIdempotency) Release(ctx context.Context, key string) error {
//     return s.client.Del(ctx, key).Err()
// }

// // WithIdempotency — usecase helper, wraps operation với idempotency check
// func WithIdempotency[T any](
//     ctx context.Context,
//     store IdempotencyStore,
//     key string,
//     ttl time.Duration,
//     fn func(ctx context.Context) (T, error),
// ) (T, error) {
//     var zero T

//     if key == "" { return fn(ctx) } // no idempotency key → execute normally

//     // Check cached result
//     var cached T
//     found, err := store.GetResult(ctx, key, &cached)
//     if err != nil {
//         if errors.Is(err, ErrInProgress) {
//             return zero, apperr.New(apperr.KindBusiness, "CONFLICT",
//                 "concurrent request with same idempotency key")
//         }
//         return zero, apperr.Internal("idempotency.check", err)
//     }
//     if found { return cached, nil } // return cached result, không execute lại

//     // Acquire lock
//     acquired, err := store.Acquire(ctx, key, ttl)
//     if err != nil || !acquired {
//         return zero, apperr.New(apperr.KindBusiness, "CONFLICT",
//             "idempotency key conflict")
//     }

//     // Execute
//     result, execErr := fn(ctx)
//     if execErr != nil {
//         store.Release(ctx, key) // allow retry
//         return zero, execErr
//     }

//     // Cache result
//     store.SetResult(ctx, key, result, ttl) // ignore error — best effort
//     return result, nil
// }
