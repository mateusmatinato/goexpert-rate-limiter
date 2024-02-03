package ratelimiter

import "errors"

var (
	ErrBlocked = errors.New("you have reached the maximum number of requests or actions " +
		"allowed within a certain time frame")
	ErrInvalidToken      = errors.New("api_key header is invalid")
	ErrInternal          = errors.New("error on rate-limiter")
	ErrNotBlocking       = errors.New("you must block by ip, token or both")
	ErrTokenListEmpty    = errors.New("you must provide a list of tokens when blocking by token")
	ErrInvalidLimitIP    = errors.New("you must provide a limit greater than 0 when blocking by ip")
	ErrInvalidLimitToken = errors.New("you must provide a limit greater than 0 when blocking by token")
	ErrGettingIP         = errors.New("error getting ip")
	ErrInvalidSetup      = errors.New("you need to activate the limit by ip to allow calls without token")
)
