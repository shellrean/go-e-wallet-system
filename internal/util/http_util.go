package util

import (
	"errors"
	"shellrean.id/belajar-auth/domain"
)

func GetHttpStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrAuthFailed):
		return 401
	case errors.Is(err, domain.ErrUsernameTaken):
		return 400
	case errors.Is(err, domain.ErrOtpInvalid):
		return 400
	case errors.Is(err, domain.ErrPinInvalid):
		return 400
	default:
		return 500
	}
}
