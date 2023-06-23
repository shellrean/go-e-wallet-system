package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"shellrean.id/belajar-auth/domain"
	"shellrean.id/belajar-auth/dto"
)

type topUpService struct {
	notificationService domain.NotificationService
	midtransService     domain.MidtransService
	topUpRepository     domain.TopUpRepository
	accountRepository   domain.AccountRepository
}

func NewTopUp(notificationService domain.NotificationService,
	midtransService domain.MidtransService,
	topUpRepository domain.TopUpRepository,
	accountRepository domain.AccountRepository) domain.TopUpService {
	return &topUpService{
		notificationService: notificationService,
		midtransService:     midtransService,
		topUpRepository:     topUpRepository,
		accountRepository:   accountRepository,
	}
}

func (t topUpService) InitializeTopUp(ctx context.Context, req dto.TopUpReq) (dto.TopUpRes, error) {
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
	topUp, err := t.topUpRepository.FindById(ctx, id)
	if err != nil {
		return err
	}

	if topUp == (domain.TopUp{}) {
		return errors.New("top-up request not found")
	}

	account, err := t.accountRepository.FindByUserID(ctx, topUp.UserID)
	if err != nil {
		return err
	}
	if account == (domain.Account{}) {
		return domain.ErrAccountNotFound
	}

	account.Balance += topUp.Amount
	err = t.accountRepository.Update(ctx, &account)
	if err != nil {
		return err
	}

	data := map[string]string{
		"amount": fmt.Sprintf("%.2f", topUp.Amount),
	}

	_ = t.notificationService.Insert(ctx, account.UserId, "TOPUP_SUCCESS", data)

	return err
}
