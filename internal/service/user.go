package service

import (
	"context"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"shellrean.id/belajar-auth/domain"
	"shellrean.id/belajar-auth/dto"
	"shellrean.id/belajar-auth/internal/util"
	"time"
)

type userService struct {
	userRepository    domain.UserRepository
	accountRepository domain.AccountRepository
	factorRepository  domain.FactorRepository
	cacheRepository   domain.CacheRepository
	emailService      domain.EmailService
}

func NewUser(userRepository domain.UserRepository,
	accountRepository domain.AccountRepository,
	factorRepository domain.FactorRepository,
	cacheRepository domain.CacheRepository,
	emailService domain.EmailService) domain.UserService {
	return &userService{
		userRepository:    userRepository,
		accountRepository: accountRepository,
		factorRepository:  factorRepository,
		cacheRepository:   cacheRepository,
		emailService:      emailService,
	}
}

func (u userService) Authenticate(ctx context.Context, req dto.AuthReq) (dto.AuthRes, error) {
	user, err := u.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return dto.AuthRes{}, err
	}
	if user == (domain.User{}) {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	if !user.EmailVerifiedAtDB.Valid {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	token := util.GenerateRandomString(32)

	userData := dto.UserData{
		ID:       user.ID,
		FullName: user.FullName,
		Phone:    user.Phone,
		Username: user.Username,
	}
	userJson, err := json.Marshal(userData)
	if err != nil {
		return dto.AuthRes{}, err
	}
	if err := u.cacheRepository.Set("user:"+token, userJson); err != nil {
		return dto.AuthRes{}, err
	}

	return dto.AuthRes{
		Token: token,
	}, nil
}

func (u userService) ValidateToken(ctx context.Context, token string) (dto.UserData, error) {
	data, err := u.cacheRepository.Get("user:" + token)
	if err != nil {
		return dto.UserData{}, domain.ErrAuthFailed
	}
	var user dto.UserData
	if err := json.Unmarshal(data, &user); err != nil || user.ID == 0 {
		return dto.UserData{}, domain.ErrAuthFailed
	}

	return user, nil
}

func (u userService) Register(ctx context.Context, req dto.UserRegisterReq) (dto.UserRegisterRes, error) {
	exist, err := u.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return dto.UserRegisterRes{}, err
	}

	if exist != (domain.User{}) {
		return dto.UserRegisterRes{}, domain.ErrUsernameTaken
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return dto.UserRegisterRes{}, err
	}

	user := domain.User{
		FullName: req.FullName,
		Phone:    req.Phone,
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPass),
	}

	err = domain.RunInTransaction(ctx, u.userRepository, func(txCtx context.Context) error {
		if err := u.userRepository.Insert(txCtx, &user); err != nil {
			return err
		}

		if u.accountRepository != nil {
			account := domain.Account{
				UserId:        user.ID,
				AccountNumber: util.GenerateRandomNumber(10),
				Balance:       0,
			}
			if err := u.accountRepository.Insert(txCtx, &account); err != nil {
				return err
			}
		}

		if u.factorRepository != nil && req.PIN != "" {
			hashedPin, err := bcrypt.GenerateFromPassword([]byte(req.PIN), 12)
			if err != nil {
				return err
			}
			factor := domain.Factor{
				UserID: user.ID,
				PIN:    string(hashedPin),
			}
			if err := u.factorRepository.Insert(txCtx, &factor); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Printf("error: %s", err.Error())
		return dto.UserRegisterRes{}, err
	}

	otpCode := util.GenerateRandomNumber(4)
	referenceId := util.GenerateRandomString(16)

	if err := u.cacheRepository.Set("otp:"+referenceId, []byte(otpCode)); err != nil {
		return dto.UserRegisterRes{}, err
	}
	if err := u.cacheRepository.Set("user-ref:"+referenceId, []byte(user.Username)); err != nil {
		return dto.UserRegisterRes{}, err
	}

	if err := u.emailService.Send(req.Email, "OTP Code", "otp anda "+otpCode); err != nil {
		log.Printf("error sending otp email: %s", err.Error())
		return dto.UserRegisterRes{}, err
	}

	return dto.UserRegisterRes{
		ReferenceID: referenceId,
	}, nil
}

func (u userService) ValidateOTP(ctx context.Context, req dto.ValidateOtpReq) error {
	val, err := u.cacheRepository.Get("otp:" + req.ReferenceID)
	if err != nil {
		return domain.ErrOtpInvalid
	}
	otp := string(val)
	if otp != req.OTP {
		return domain.ErrOtpInvalid
	}
	val, err = u.cacheRepository.Get("user-ref:" + req.ReferenceID)
	if err != nil {
		return domain.ErrOtpInvalid
	}
	user, err := u.userRepository.FindByUsername(ctx, string(val))
	if err != nil {
		return err
	}
	if user == (domain.User{}) {
		return domain.ErrOtpInvalid
	}

	user.EmailVerifiedAt = time.Now()
	if err := u.userRepository.Update(ctx, &user); err != nil {
		return err
	}

	_ = u.cacheRepository.Delete("otp:" + req.ReferenceID)
	_ = u.cacheRepository.Delete("user-ref:" + req.ReferenceID)
	return nil
}
