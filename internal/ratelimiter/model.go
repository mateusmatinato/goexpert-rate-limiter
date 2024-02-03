package ratelimiter

import (
	"fmt"
	"time"
)

type LimitOrigin string

const (
	OriginToken LimitOrigin = "token"
	OriginIP    LimitOrigin = "ip"
)

func generateKey(limitOrigin LimitOrigin, value string) string {
	return fmt.Sprintf("%s:%s", limitOrigin, value)
}

type Params struct {
	BlockByIP       bool
	BlockByToken    bool
	BlockTimeToken  time.Duration
	BlockTimeIP     time.Duration
	LimitIPByMinute int
	TokenList       map[string]TokenInfo
}

type TokenInfo struct {
	ID                  string
	MaxRequestsByMinute int
}
