package service

import (
	"context"
	"encoding/json"
	"fmt"
	"shellrean.id/belajar-auth/domain"
	"shellrean.id/belajar-auth/dto"
	"shellrean.id/belajar-auth/internal/util"
	"time"
)

type transactionService struct {
	accountRepository     domain.AccountRepository
	transactionRepository domain.TransactionRepository
	cacheRepository       domain.CacheRepository
	notificationService   domain.NotificationService
}

func NewTransaction(accountRepository domain.AccountRepository,
	transactionRepository domain.TransactionRepository,
	cacheRepository domain.CacheRepository,
	notificationService domain.NotificationService) domain.TransactionService {
	return &transactionService{
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		cacheRepository:       cacheRepository,
		notificationService:   notificationService,
	}
}

func (t transactionService) TransferInquiry(ctx context.Context, req dto.TransferInquiryReq) (dto.TransferInquiryRes, error) {
	user, ok := ctx.Value("x-user").(dto.UserData)
	if !ok {
		return dto.TransferInquiryRes{}, domain.ErrAuthFailed
	}

	req.Amount = util.RoundCurrency(req.Amount)
	if req.Amount <= 0 {
		return dto.TransferInquiryRes{}, domain.ErrInvalidAmount
	}

	myAccount, err := t.accountRepository.FindByUserID(ctx, user.ID)
	if err != nil {
		return dto.TransferInquiryRes{}, err
	}

	if myAccount == (domain.Account{}) {
		return dto.TransferInquiryRes{}, domain.ErrAccountNotFound
	}

	dofAccount, err := t.accountRepository.FindByAccountNumber(ctx, req.AccountNumber)
	if err != nil {
		return dto.TransferInquiryRes{}, err
	}

	if dofAccount == (domain.Account{}) {
		return dto.TransferInquiryRes{}, domain.ErrAccountNotFound
	}

	if myAccount.AccountNumber == dofAccount.AccountNumber {
		return dto.TransferInquiryRes{}, domain.ErrSelfTransfer
	}

	if myAccount.Balance < req.Amount {
		return dto.TransferInquiryRes{}, domain.ErrInsufficientBalance
	}

	inquiryKey := util.GenerateRandomString(32)
	req.UserID = user.ID

	jsonData, err := json.Marshal(req)
	if err != nil {
		return dto.TransferInquiryRes{}, err
	}
	if err := t.cacheRepository.Set(inquiryKey, jsonData); err != nil {
		return dto.TransferInquiryRes{}, err
	}

	return dto.TransferInquiryRes{
		InquiryKey: inquiryKey,
	}, nil
}

func (t transactionService) TransferExecute(ctx context.Context, req dto.TransferExecuteReq) error {
	user, ok := ctx.Value("x-user").(dto.UserData)
	if !ok {
		return domain.ErrAuthFailed
	}

	val, err := t.cacheRepository.Get(req.InquiryKey)
	if err != nil {
		return domain.ErrInquiryNotFound
	}

	var reqInq dto.TransferInquiryReq
	if err := json.Unmarshal(val, &reqInq); err != nil {
		return domain.ErrInquiryNotFound
	}
	if reqInq == (dto.TransferInquiryReq{}) {
		return domain.ErrInquiryNotFound
	}

	reqInq.Amount = util.RoundCurrency(reqInq.Amount)
	if reqInq.Amount <= 0 {
		return domain.ErrInvalidAmount
	}

	if reqInq.UserID != 0 && reqInq.UserID != user.ID {
		return domain.ErrInquiryNotFound
	}

	_ = t.cacheRepository.Delete(req.InquiryKey)

	var myAccount, dofAccount domain.Account
	err = domain.RunInTransaction(ctx, t.accountRepository, func(txCtx context.Context) error {
		var err error
		myAccount, err = t.accountRepository.FindByUserID(txCtx, user.ID)
		if err != nil {
			return err
		}
		if myAccount == (domain.Account{}) {
			return domain.ErrAccountNotFound
		}

		dofAccount, err = t.accountRepository.FindByAccountNumber(txCtx, reqInq.AccountNumber)
		if err != nil {
			return err
		}
		if dofAccount == (domain.Account{}) {
			return domain.ErrAccountNotFound
		}

		if myAccount.AccountNumber == dofAccount.AccountNumber {
			return domain.ErrSelfTransfer
		}

		if myAccount.Balance < reqInq.Amount {
			return domain.ErrInsufficientBalance
		}

		now := time.Now()
		debitTransaction := domain.Transaction{
			AccountId:           myAccount.ID,
			SofNumber:           myAccount.AccountNumber,
			DofNumber:           dofAccount.AccountNumber,
			TransactionType:     "D",
			Amount:              reqInq.Amount,
			TransactionDatetime: now,
		}
		if err := t.transactionRepository.Insert(txCtx, &debitTransaction); err != nil {
			return err
		}

		creditTransaction := domain.Transaction{
			AccountId:           dofAccount.ID,
			SofNumber:           myAccount.AccountNumber,
			DofNumber:           dofAccount.AccountNumber,
			TransactionType:     "C",
			Amount:              reqInq.Amount,
			TransactionDatetime: now,
		}
		if err := t.transactionRepository.Insert(txCtx, &creditTransaction); err != nil {
			return err
		}

		myAccount.Balance = util.RoundCurrency(myAccount.Balance - reqInq.Amount)
		if err := t.accountRepository.Update(txCtx, &myAccount); err != nil {
			return err
		}

		dofAccount.Balance = util.RoundCurrency(dofAccount.Balance + reqInq.Amount)
		if err := t.accountRepository.Update(txCtx, &dofAccount); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	go t.notificationAfterTransfer(myAccount, dofAccount, reqInq.Amount)
	return nil
}

func (t transactionService) notificationAfterTransfer(sofAccount domain.Account, dofAccount domain.Account, amount float64) {
	data := map[string]string{
		"amount": fmt.Sprintf("%.2f", amount),
	}

	_ = t.notificationService.Insert(context.Background(), sofAccount.UserId, "TRANSFER", data)
	_ = t.notificationService.Insert(context.Background(), dofAccount.UserId, "TRANSFER_DES", data)
}
