package domain

import "errors"

var ErrAuthFailed = errors.New("error authentication failed")
var ErrUsernameTaken = errors.New("username already taken")
var ErrOtpInvalid = errors.New("otp invalid")
