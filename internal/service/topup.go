package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"shellrean.id/belajar-auth/domain"
	"shellrean.id/belajar-auth/dto"
	"shellrean.id/belajar-auth/internal/util"
	"time"
)

type topUpService struct {
	notificationService   domain.NotificationService
	midtransService       domain.MidtransService
	topUpRepository       domain.TopUpRepository
	accountRepository     domain.AccountRepository
	transactionRepository domain.TransactionRepository
}

func NewTopUp(notificationService domain.NotificationService,
	midtransService domain.MidtransService,
	topUpRepository domain.TopUpRepository,
	accountRepository domain.AccountRepository,
	transactionRepository domain.TransactionRepository) domain.TopUpService {
	return &topUpService{
		notificationService:   notificationService,
		midtransService:       midtransService,
		topUpRepository:       topUpRepository,
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
	}
}

func (t topUpService) InitializeTopUp(ctx context.Context, req dto.TopUpReq) (dto.TopUpRes, error) {
	req.Amount = util.RoundCurrency(req.Amount)
	if req.Amount <= 0 {
		return dto.TopUpRes{}, domain.ErrInvalidAmount
	}

	topUp := domain.TopUp{
		ID:     uuid.NewString(),
		UserID: req.UserID,
		Status: 0,
		Amount: req.Amount,
	}
	err := t.midtransService.GenerateSnapURL(ctx, &topUp)
	if err != nil {
		return dto.TopUpRes{}, err
	}

	err = t.topUpRepository.Insert(ctx, &topUp)
	if err != nil {
		return dto.TopUpRes{}, err
	}

	return dto.TopUpRes{
		SnapURL: topUp.SnapURL,
	}, nil
}

func (t topUpService) ConfirmedTopUp(ctx context.Context, id string) error {
	var account domain.Account
	var topUp domain.TopUp

	err := domain.RunInTransaction(ctx, t.topUpRepository, func(txCtx context.Context) error {
		var err error
		topUp, err = t.topUpRepository.FindById(txCtx, id)
		if err != nil {
			return err
		}

		if topUp == (domain.TopUp{}) {
			return domain.ErrTopUpNotFound
		}

		if topUp.Status == 1 {
			return domain.ErrTopUpAlreadyConfirmed
		}

		account, err = t.accountRepository.FindByUserID(txCtx, topUp.UserID)
		if err != nil {
			return err
		}
		if account == (domain.Account{}) {
			return domain.ErrAccountNotFound
		}

		err = t.transactionRepository.Insert(txCtx, &domain.Transaction{
			AccountId:           account.ID,
			SofNumber:           "00",
			DofNumber:           account.AccountNumber,
			TransactionType:     "C",
			Amount:              topUp.Amount,
			TransactionDatetime: time.Now(),
		})
		if err != nil {
			return err
		}

		account.Balance = util.RoundCurrency(account.Balance + topUp.Amount)
		if err := t.accountRepository.Update(txCtx, &account); err != nil {
			return err
		}

		topUp.Status = 1
		if err := t.topUpRepository.Update(txCtx, &topUp); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	data := map[string]string{
		"amount": fmt.Sprintf("%.2f", topUp.Amount),
	}

	_ = t.notificationService.Insert(ctx, account.UserId, "TOPUP_SUCCESS", data)

	return nil
}
