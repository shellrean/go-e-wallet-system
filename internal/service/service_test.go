package service_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"shellrean.id/belajar-auth/domain"
	"shellrean.id/belajar-auth/dto"
	"shellrean.id/belajar-auth/internal/service"
	"shellrean.id/belajar-auth/internal/util"
)

type memoryCache struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func newMemoryCache() *memoryCache {
	return &memoryCache{data: make(map[string][]byte)}
}

func (m *memoryCache) Get(key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	if !ok {
		return nil, errors.New("not found")
	}
	return val, nil
}

func (m *memoryCache) Set(key string, entry []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = entry
	return nil
}

func (m *memoryCache) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}

type fakeAccountRepo struct {
	mu       sync.Mutex
	accounts map[int64]domain.Account
}

func newFakeAccountRepo() *fakeAccountRepo {
	return &fakeAccountRepo{accounts: make(map[int64]domain.Account)}
}

func (f *fakeAccountRepo) FindByUserID(ctx context.Context, id int64) (domain.Account, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, a := range f.accounts {
		if a.UserId == id {
			return a, nil
		}
	}
	return domain.Account{}, nil
}

func (f *fakeAccountRepo) FindByAccountNumber(ctx context.Context, accNumber string) (domain.Account, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, a := range f.accounts {
		if a.AccountNumber == accNumber {
			return a, nil
		}
	}
	return domain.Account{}, nil
}

func (f *fakeAccountRepo) Insert(ctx context.Context, account *domain.Account) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	account.ID = int64(len(f.accounts) + 1)
	f.accounts[account.ID] = *account
	return nil
}

func (f *fakeAccountRepo) Update(ctx context.Context, account *domain.Account) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.accounts[account.ID] = *account
	return nil
}

type fakeTransactionRepo struct {
	mu   sync.Mutex
	txns []domain.Transaction
}

func (f *fakeTransactionRepo) Insert(ctx context.Context, transaction *domain.Transaction) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	transaction.ID = int64(len(f.txns) + 1)
	f.txns = append(f.txns, *transaction)
	return nil
}

type fakeTopUpRepo struct {
	mu     sync.Mutex
	topups map[string]domain.TopUp
}

func newFakeTopUpRepo() *fakeTopUpRepo {
	return &fakeTopUpRepo{topups: make(map[string]domain.TopUp)}
}

func (f *fakeTopUpRepo) FindById(ctx context.Context, id string) (domain.TopUp, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.topups[id], nil
}

func (f *fakeTopUpRepo) Insert(ctx context.Context, t *domain.TopUp) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.topups[t.ID] = *t
	return nil
}

func (f *fakeTopUpRepo) Update(ctx context.Context, t *domain.TopUp) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.topups[t.ID] = *t
	return nil
}

type noopNotificationService struct{}

func (noopNotificationService) FindByUser(ctx context.Context, user int64) ([]dto.NotificationData, error) {
	return nil, nil
}
func (noopNotificationService) Insert(ctx context.Context, userId int64, code string, data map[string]string) error {
	return nil
}

type fakeMidtransService struct{}

func (fakeMidtransService) GenerateSnapURL(ctx context.Context, t *domain.TopUp) error {
	t.SnapURL = "https://app.sandbox.midtrans.com/snap/v2/vtweb/" + t.ID
	return nil
}
func (fakeMidtransService) VerifyPayment(ctx context.Context, orderId string) (bool, error) {
	return true, nil
}

func TestTransferReplayAndValidation(t *testing.T) {
	accRepo := newFakeAccountRepo()
	_ = accRepo.Insert(context.Background(), &domain.Account{UserId: 1, AccountNumber: "111111", Balance: 100000})
	_ = accRepo.Insert(context.Background(), &domain.Account{UserId: 2, AccountNumber: "222222", Balance: 50000})

	txRepo := &fakeTransactionRepo{}
	cache := newMemoryCache()
	svc := service.NewTransaction(accRepo, txRepo, cache, noopNotificationService{})

	ctxUser1 := context.WithValue(context.Background(), "x-user", dto.UserData{ID: 1, Username: "alice"})

	// 1. Reject negative amount
	_, err := svc.TransferInquiry(ctxUser1, dto.TransferInquiryReq{AccountNumber: "222222", Amount: -1000})
	if !errors.Is(err, domain.ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}

	// 2. Reject self-transfer
	_, err = svc.TransferInquiry(ctxUser1, dto.TransferInquiryReq{AccountNumber: "111111", Amount: 10000})
	if !errors.Is(err, domain.ErrSelfTransfer) {
		t.Fatalf("expected ErrSelfTransfer, got %v", err)
	}

	// 3. Valid inquiry & execute
	inq, err := svc.TransferInquiry(ctxUser1, dto.TransferInquiryReq{AccountNumber: "222222", Amount: 40000})
	if err != nil {
		t.Fatalf("unexpected inquiry error: %v", err)
	}

	if err := svc.TransferExecute(ctxUser1, dto.TransferExecuteReq{InquiryKey: inq.InquiryKey}); err != nil {
		t.Fatalf("unexpected execute error: %v", err)
	}

	// 4. Reject replay of the same InquiryKey
	if err := svc.TransferExecute(ctxUser1, dto.TransferExecuteReq{InquiryKey: inq.InquiryKey}); !errors.Is(err, domain.ErrInquiryNotFound) {
		t.Fatalf("expected ErrInquiryNotFound on replay, got %v", err)
	}

	acc1, _ := accRepo.FindByUserID(context.Background(), 1)
	acc2, _ := accRepo.FindByUserID(context.Background(), 2)
	if acc1.Balance != 60000 || acc2.Balance != 90000 {
		t.Fatalf("unexpected balances: acc1=%v, acc2=%v", acc1.Balance, acc2.Balance)
	}
}

func TestTopUpIdempotency(t *testing.T) {
	accRepo := newFakeAccountRepo()
	_ = accRepo.Insert(context.Background(), &domain.Account{UserId: 1, AccountNumber: "111111", Balance: 10000})

	topUpRepo := newFakeTopUpRepo()
	txRepo := &fakeTransactionRepo{}
	svc := service.NewTopUp(noopNotificationService{}, fakeMidtransService{}, topUpRepo, accRepo, txRepo)

	// Reject zero/negative top-up
	if _, err := svc.InitializeTopUp(context.Background(), dto.TopUpReq{UserID: 1, Amount: 0}); !errors.Is(err, domain.ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}

	topUpID := "order-123"
	_ = topUpRepo.Insert(context.Background(), &domain.TopUp{
		ID:     topUpID,
		UserID: 1,
		Status: 0,
		Amount: 50000,
	})

	if err := svc.ConfirmedTopUp(context.Background(), topUpID); err != nil {
		t.Fatalf("unexpected error on first confirm: %v", err)
	}

	// Second confirmation must return ErrTopUpAlreadyConfirmed and not double-credit
	if err := svc.ConfirmedTopUp(context.Background(), topUpID); !errors.Is(err, domain.ErrTopUpAlreadyConfirmed) {
		t.Fatalf("expected ErrTopUpAlreadyConfirmed, got %v", err)
	}

	acc, _ := accRepo.FindByUserID(context.Background(), 1)
	if acc.Balance != 60000 {
		t.Fatalf("expected balance 60000, got %v", acc.Balance)
	}
}

func TestHubConcurrency(t *testing.T) {
	hub := dto.NewHub()
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			ch := hub.Register(id)
			hub.Publish(id, dto.NotificationData{ID: id, Title: "Test"})
			hub.Unregister(id, ch)
		}(int64(i % 5))
	}
	wg.Wait()

	if status := util.GetHttpStatus(domain.ErrInsufficientBalance); status != 400 {
		t.Fatalf("expected 400 for ErrInsufficientBalance, got %d", status)
	}
	if status := util.GetHttpStatus(domain.ErrAccountNotFound); status != 404 {
		t.Fatalf("expected 404 for ErrAccountNotFound, got %d", status)
	}
}
