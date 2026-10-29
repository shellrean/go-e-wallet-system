package domain

import "errors"

var ErrAuthFailed = errors.New("error authentication failed")
var ErrUsernameTaken = errors.New("username already taken")
var ErrOtpInvalid = errors.New("otp invalid")
var ErrAccountNotFound = errors.New("account not found")
var ErrInquiryNotFound = errors.New("inquiry not found")
var ErrInsufficientBalance = errors.New("insufficient balance")
var ErrPinInvalid = errors.New("pin invalid")
var ErrInvalidAmount = errors.New("amount must be greater than zero")
var ErrSelfTransfer = errors.New("cannot transfer to the same account")
var ErrTopUpNotFound = errors.New("top-up request not found")
var ErrTopUpAlreadyConfirmed = errors.New("top-up already confirmed")
var ErrTemplateNotFound = errors.New("template not found")
