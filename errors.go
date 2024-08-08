package main

import "errors"

var (
	ErrInsufficientFundsNoReorg = errors.New("Insufficient funds; no re-org needed.")
	ErrSignalTooWeak            = errors.New("Insufficient funds; signal too weak.")
	ErrInsufficientFunds        = errors.New("Insufficient funds.")
	ErrorRecordTooLarge         = errors.New("Record too large.")
	ErrorInvalidAddress         = errors.New("Invalid Bitcoin address.")
	ErrorNeedMoreSats           = errors.New("Signal needs more sats.")
	ErrorInvalidSignature       = errors.New("Invalid signature.")
)
