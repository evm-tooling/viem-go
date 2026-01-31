package accounts

import "errors"

var (
	// ErrInvalidAddress is returned when an address is invalid.
	ErrInvalidAddress = errors.New("invalid address")
	// ErrInvalidPrivateKey is returned when a private key is invalid.
	ErrInvalidPrivateKey = errors.New("invalid private key")
	// ErrInvalidMnemonic is returned when a mnemonic is invalid.
	ErrInvalidMnemonic = errors.New("invalid mnemonic")
	// ErrInvalidHDPath is returned when an HD derivation path is invalid.
	ErrInvalidHDPath = errors.New("invalid HD derivation path")
	// ErrSigningNotSupported is returned when signing is not supported.
	ErrSigningNotSupported = errors.New("signing not supported for this account type")
	// ErrInvalidWordlist is returned when a wordlist is invalid.
	ErrInvalidWordlist = errors.New("invalid wordlist")
)
