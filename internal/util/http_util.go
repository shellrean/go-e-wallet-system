package util

import (
	"errors"
	"shellrean.id/belajar-auth/domain"
)

func GetHttpStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrAuthFailed):
		return 401
	case errors.Is(err, domain.ErrAccountNotFound),
		errors.Is(err, domain.ErrInquiryNotFound),
		errors.Is(err, domain.ErrTopUpNotFound),
		errors.Is(err, domain.ErrTemplateNotFound):
		return 404
	case errors.Is(err, domain.ErrUsernameTaken),
		errors.Is(err, domain.ErrOtpInvalid),
		errors.Is(err, domain.ErrPinInvalid),
		errors.Is(err, domain.ErrInsufficientBalance),
		errors.Is(err, domain.ErrInvalidAmount),
		errors.Is(err, domain.ErrSelfTransfer),
		errors.Is(err, domain.ErrTopUpAlreadyConfirmed):
		return 400
	default:
		return 500
	}
}
