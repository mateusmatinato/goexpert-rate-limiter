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
	LimitByIP       bool
	LimitByToken    bool
	BlockTimeToken  time.Duration
	BlockTimeIP     time.Duration
	LimitIPBySecond int
	TokenList       map[string]int
}
