# II.5 SQL Transaction - Commit và Rollback

## 🎯 Tổng Quan

SQL Transaction trong dự án BDSPro Microservices được quản lý thông qua GORM với pattern interface và implementation. Hệ thống cung cấp các phương thức để bắt đầu, commit, rollback transaction và đảm bảo tính toàn vẹn dữ liệu.

---

## 📋 Interface Transaction

### **ITransaction Interface**
```go
type ITransaction interface {
    WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
    StartTransaction(ctx context.Context) context.Context
    RollbackTransaction(ctx context.Context) error
    CommitTransaction(ctx context.Context) error
}
```

### **Các Method Chính:**
- **`WithTransaction`**: Thực hiện function trong transaction, tự động commit/rollback
- **`StartTransaction`**: Bắt đầu transaction thủ công
- **`RollbackTransaction`**: Rollback transaction
- **`CommitTransaction`**: Commit transaction

---

## 🔧 Implementation Transaction

### **TransactionGorm Implementation**
```go
// @bind: service/internal/interface.ITransaction
type TransactionGorm struct {
    DB *gorm.DB
}

func NewTransactionGorm(DB *gorm.DB) *TransactionGorm {
    return &TransactionGorm{DB: DB}
}

// Tự động quản lý transaction
func (t *TransactionGorm) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
    return t.DB.Transaction(func(tx *gorm.DB) error {
        return fn(context.WithValue(ctx, "tx", tx))
    })
}

// Bắt đầu transaction thủ công
func (t *TransactionGorm) StartTransaction(ctx context.Context) context.Context {
    tx := t.DB.Begin()
    return context.WithValue(ctx, "tx", tx)
}

// Rollback transaction
func (t *TransactionGorm) RollbackTransaction(ctx context.Context) error {
    return GetDB(ctx, t.DB).Rollback().Error
}

// Commit transaction
func (t *TransactionGorm) CommitTransaction(ctx context.Context) error {
    return GetDB(ctx, t.DB).Commit().Error
}

// Helper function để lấy DB instance
func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
    if tx, ok := ctx.Value("tx").(*gorm.DB); ok && tx != nil {
        return tx.WithContext(ctx)
    }
    return defaultDB
}
```

---

## 🚀 Cách Sử Dụng Transaction

### **1. Sử Dụng WithTransaction (Khuyến nghị)**

#### **Trong Usecase Layer:**
```go
func (u *userUsecase) CreateUserWithProfile(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
    var result *dto.CreateUserResponse
    
    err := u.Transaction.WithTransaction(ctx, func(ctx context.Context) error {
        // Tạo user
        user, err := u.userRepo.Create(ctx, &entity.User{
            Name:  req.Name,
            Email: req.Email,
        })
        if err != nil {
            return err // Tự động rollback
        }

        // Tạo profile
        profile, err := u.profileRepo.Create(ctx, &entity.Profile{
            UserID: user.ID,
            Bio:    req.Bio,
        })
        if err != nil {
            return err // Tự động rollback
        }

        result = &dto.CreateUserResponse{
            User:    user,
            Profile: profile,
        }
        
        return nil // Tự động commit
    })
    
    if err != nil {
        return nil, err
    }
    
    return result, nil
}
```

#### **Với User Context:**
```go
func (t *PostgreTransaction) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
    return t.db.Transaction(func(tx *gorm.DB) error {
        // Truyền tx và user info vào context
        profileId := _utils.GetProfileIdWithContext(ctx)
        organizationId := _utils.GetOrganizationIdFromContext(ctx)

        newCtx := context.WithValue(ctx, "tx", tx)
        newCtx = context.WithValue(newCtx, "profileId", profileId)
        newCtx = context.WithValue(newCtx, "organizationId", organizationId)

        return fn(newCtx)
    })
}
```

### **2. Sử Dụng Manual Transaction**

#### **Trong Usecase Layer:**
```go
func (p *paymentUsecase) MakePayment(ctx context.Context, req *paymentpb.PaymentRequest) (*paymentpb.WalletTransaction, error) {
    // Bắt đầu transaction
    tx := p.walletRepository.BeginTx(ctx)
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }()

    // Tạo transaction context
    txCtx := context.WithValue(ctx, "tx", tx)

    // Get wallet
    wallet, err := p.walletRepository.GetWalletByUserId(txCtx, userID)
    if err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to get wallet: %w", err)
    }

    // Check balance
    if wallet.Balance < req.Amount {
        tx.Rollback()
        return nil, fmt.Errorf("insufficient balance")
    }

    // Create transaction
    transaction := &entity.WalletTransaction{
        WalletId: wallet.Id,
        Type:     "PAYMENT",
        Amount:   -req.Amount,
        Status:   "COMPLETED",
    }

    err = p.walletTransactionRepository.CreateTransaction(txCtx, transaction)
    if err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to create transaction: %w", err)
    }

    // Update wallet balance
    wallet.Balance -= req.Amount
    err = p.walletRepository.UpdateWallet(txCtx, wallet)
    if err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to update wallet balance: %w", err)
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }

    return transaction, nil
}
```

### **3. Sử Dụng Start/Commit/Rollback**

#### **Trong Usecase Layer:**
```go
func (u *userUsecase) ComplexUserOperation(ctx context.Context, req *dto.ComplexRequest) error {
    // Bắt đầu transaction
    txCtx := u.Transaction.StartTransaction(ctx)
    
    // Thực hiện các operations
    user, err := u.userRepo.Create(txCtx, &entity.User{...})
    if err != nil {
        u.Transaction.RollbackTransaction(txCtx)
        return err
    }

    profile, err := u.profileRepo.Create(txCtx, &entity.Profile{...})
    if err != nil {
        u.Transaction.RollbackTransaction(txCtx)
        return err
    }

    // Commit transaction
    if err := u.Transaction.CommitTransaction(txCtx); err != nil {
        return err
    }

    return nil
}
```

---

## 🏗️ Repository Pattern với Transaction

### **Repository Implementation:**
```go
// infra/postgre/user_postgres.go
func (r *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
    // Sử dụng GetDB để lấy đúng DB instance (transaction hoặc default)
    db := GetDB(ctx, r.db)
    
    err := db.WithContext(ctx).Create(user).Error
    if err != nil {
        return nil, err
    }
    
    return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
    db := GetDB(ctx, r.db)
    
    err := db.WithContext(ctx).Save(user).Error
    if err != nil {
        return nil, err
    }
    
    return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id uint64) error {
    db := GetDB(ctx, r.db)
    
    return db.WithContext(ctx).Delete(&entity.User{}, id).Error
}
```

---

## 🔄 Transaction Context Flow

### **Context Propagation:**
```
1. Usecase gọi WithTransaction
2. Transaction bắt đầu và lưu tx vào context
3. Context được truyền xuống Repository
4. Repository sử dụng GetDB để lấy tx từ context
5. Tất cả operations sử dụng cùng transaction
6. Tự động commit/rollback khi function kết thúc
```

### **Context Keys:**
- **`"tx"`**: Chứa GORM transaction instance
- **`"profileId"`**: User profile ID (nếu có)
- **`"organizationId"`**: Organization ID (nếu có)

---

## ⚠️ Lưu Ý Quan Trọng

### **1. Quy Tắc Sử Dụng**
- **Ưu tiên `WithTransaction`** cho hầu hết trường hợp
- **Sử dụng manual transaction** khi cần kiểm soát chi tiết
- **Luôn sử dụng `GetDB(ctx, defaultDB)`** trong repository
- **Không mix transaction và non-transaction operations**

### **2. Error Handling**
```go
// ✅ ĐÚNG - Return error để WithTransaction tự động rollback
err := u.Transaction.WithTransaction(ctx, func(ctx context.Context) error {
    user, err := u.userRepo.Create(ctx, user)
    if err != nil {
        return err // Tự động rollback
    }
    
    profile, err := u.profileRepo.Create(ctx, profile)
    if err != nil {
        return err // Tự động rollback
    }
    
    return nil // Tự động commit
})

// ❌ SAI - Không return error
err := u.Transaction.WithTransaction(ctx, func(ctx context.Context) error {
    user, err := u.userRepo.Create(ctx, user)
    if err != nil {
        // Không return error - transaction vẫn commit
        log.Printf("Error: %v", err)
    }
    return nil
})
```

### **3. Performance Considerations**
- **Transaction scope nhỏ nhất có thể**
- **Tránh long-running transactions**
- **Sử dụng proper indexing**
- **Monitor deadlocks**

### **4. Best Practices**
- **Luôn validate input trước khi bắt đầu transaction**
- **Sử dụng defer để cleanup resources**
- **Log transaction operations cho debugging**
- **Test transaction rollback scenarios**

---

## 🔍 Debugging Transaction

### **1. Transaction Logging**
```go
func (t *TransactionGorm) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
    log.Printf("Starting transaction")
    
    err := t.DB.Transaction(func(tx *gorm.DB) error {
        log.Printf("Transaction started")
        
        err := fn(context.WithValue(ctx, "tx", tx))
        if err != nil {
            log.Printf("Transaction failed: %v", err)
        } else {
            log.Printf("Transaction succeeded")
        }
        
        return err
    })
    
    if err != nil {
        log.Printf("Transaction rolled back: %v", err)
    } else {
        log.Printf("Transaction committed")
    }
    
    return err
}
```

### **2. Context Debugging**
```go
func debugTransactionContext(ctx context.Context) {
    if tx, ok := ctx.Value("tx").(*gorm.DB); ok {
        log.Printf("Transaction context found: %p", tx)
    } else {
        log.Printf("No transaction context")
    }
}
```

---

## 🎯 Tóm Tắt

| Method | Mục đích | Khi nào sử dụng |
|--------|----------|-----------------|
| **`WithTransaction`** | Tự động quản lý transaction | Hầu hết trường hợp |
| **`StartTransaction`** | Bắt đầu transaction thủ công | Khi cần kiểm soát chi tiết |
| **`CommitTransaction`** | Commit transaction | Với manual transaction |
| **`RollbackTransaction`** | Rollback transaction | Với manual transaction |
| **`GetDB`** | Lấy DB instance từ context | Trong repository layer |

**Quy tắc vàng**: **Luôn sử dụng `WithTransaction` cho hầu hết trường hợp và `GetDB(ctx, defaultDB)` trong repository để đảm bảo transaction consistency!**
