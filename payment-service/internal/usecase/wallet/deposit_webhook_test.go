package wallet

import (
	"context"
	"math"
	"sync"
	"testing"

	walletdomain "payment/internal/domain/wallet"
)

type depositTxKey struct{}
type depositTx struct {
	unlock func()
	once   sync.Once
}

func (t *depositTx) Commit() error { return t.Rollback() }
func (t *depositTx) Rollback() error {
	t.once.Do(func() {
		if t.unlock != nil {
			t.unlock()
		}
	})
	return nil
}
func (t *depositTx) Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, depositTxKey{}, t)
}

type depositWalletStore struct {
	WalletRepository
	wallet  walletdomain.Wallet
	updates int
}

func (s *depositWalletStore) BeginTx(context.Context) Transaction { return &depositTx{} }
func (s *depositWalletStore) GetWalletByIdForUpdate(context.Context, uint32) (*walletdomain.Wallet, error) {
	v := s.wallet
	return &v, nil
}
func (s *depositWalletStore) UpdateWallet(_ context.Context, w *walletdomain.Wallet) error {
	s.wallet = *w
	s.updates++
	return nil
}

type depositTransactionStore struct {
	WalletTransactionRepository
	mu    sync.Mutex
	value walletdomain.WalletTransaction
}

func (s *depositTransactionStore) GetTransactionByCodeForUpdate(ctx context.Context, _ string) (*walletdomain.WalletTransaction, error) {
	s.mu.Lock()
	ctx.Value(depositTxKey{}).(*depositTx).unlock = s.mu.Unlock
	v := s.value
	return &v, nil
}
func (s *depositTransactionStore) UpdateTransaction(_ context.Context, v *walletdomain.WalletTransaction) error {
	s.value = *v
	return nil
}

type depositAudit struct{ count int }

func (s *depositAudit) CreateAuditLog(context.Context, *walletdomain.WalletAuditLog) error {
	s.count++
	return nil
}

func TestDepositConcurrentWebhookReplayCreditsOnce(t *testing.T) {
	wallets := &depositWalletStore{wallet: walletdomain.Wallet{Id: 7, Balance: 2000}}
	transactions := &depositTransactionStore{value: walletdomain.WalletTransaction{Id: 1, WalletId: 7, Amount: 5000, Type: walletdomain.TransactionTypeDeposit, Status: walletdomain.TransactionStatusPending}}
	audit := &depositAudit{}
	service := &paymentUsecase{walletRepository: wallets, walletTransactionRepository: transactions, walletAuditLogRepository: audit}
	results := make(chan error, 2)
	for range 2 {
		go func() { _, err := service.HandleDepositWebhook(context.Background(), "DEP123", 5000); results <- err }()
	}
	for range 2 {
		if err := <-results; err != nil {
			t.Fatal(err)
		}
	}
	if wallets.wallet.Balance != 7000 || wallets.updates != 1 || audit.count != 1 {
		t.Fatalf("balance=%d updates=%d audits=%d", wallets.wallet.Balance, wallets.updates, audit.count)
	}
}

func TestDepositRejectsInvalidLifecycleAndOverflow(t *testing.T) {
	for _, tc := range []struct {
		name    string
		kind    walletdomain.TransactionType
		status  walletdomain.TransactionStatus
		balance int64
	}{
		{"canceled", walletdomain.TransactionTypeDeposit, walletdomain.TransactionStatusCanceled, 0},
		{"wrong kind", walletdomain.TransactionTypePayment, walletdomain.TransactionStatusPending, 0},
		{"overflow", walletdomain.TransactionTypeDeposit, walletdomain.TransactionStatusPending, math.MaxInt64},
	} {
		t.Run(tc.name, func(t *testing.T) {
			wallets := &depositWalletStore{wallet: walletdomain.Wallet{Id: 7, Balance: tc.balance}}
			transactions := &depositTransactionStore{value: walletdomain.WalletTransaction{Id: 1, WalletId: 7, Amount: 5000, Type: tc.kind, Status: tc.status}}
			audit := &depositAudit{}
			service := &paymentUsecase{walletRepository: wallets, walletTransactionRepository: transactions, walletAuditLogRepository: audit}
			if _, err := service.HandleDepositWebhook(context.Background(), "DEP123", 5000); err == nil {
				t.Fatal("invalid deposit accepted")
			}
			if wallets.updates != 0 || audit.count != 0 {
				t.Fatal("rejected deposit changed balance or audit")
			}
		})
	}
}
