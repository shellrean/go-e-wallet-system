package domain

import "errors"

var ErrAuthFailed = errors.New("error authentication failed")
var ErrUsernameTaken = errors.New("username already taken")
var ErrOtpInvalid = errors.New("otp invalid")
var ErrAccountNotFound = errors.New("account not found")
var ErrInquiryNotFound = errors.New("inquiry not found")
var ErrInsufficientBalance = errors.New("insufficient balance")
