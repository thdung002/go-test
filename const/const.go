package _const

import "errors"

var (
	ERR_NONCE               = errors.New("Nonce check error")
	ERR_INSUFFICIENT        = errors.New("Insufficient amount")
	ERR_ADDRESS_LENGTH      = errors.New("Address length error")
	ERR_INVALID_PRIVATE_KEY = errors.New("Invalid private key")
	ERR_USER_NOT_FOUND      = errors.New("Not found")
	ERR_VERIFY_SIGNATURE    = errors.New("Not found")
)
