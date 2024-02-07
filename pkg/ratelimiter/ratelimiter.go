package ratelimiter

import (
	"net/http"
	"time"

	httpInternal "github.com/mateusmatinato/goexpert-rate-limiter/internal/platform/http"
	"github.com/mateusmatinato/goexpert-rate-limiter/internal/platform/redis"
	"github.com/mateusmatinato/goexpert-rate-limiter/internal/ratelimiter"
	"github.com/mateusmatinato/goexpert-rate-limiter/internal/ratelimiter/access"
	"github.com/mateusmatinato/goexpert-rate-limiter/internal/ratelimiter/blocked"
)

type (
	ParamsOptions func(params *ratelimiter.Params)
)

type TokenInfo map[string]int

type DatabaseConfig struct {
	Addr     string
	Port     int
	Password string
}

func defaultParams() ratelimiter.Params {
	return ratelimiter.Params{
		LimitByIP:       false,
		LimitByToken:    false,
		BlockTimeToken:  1 * time.Minute,
		BlockTimeIP:     1 * time.Minute,
		LimitIPBySecond: 5,
		TokenList:       make(map[string]int),
	}
}

func WithLimitByIP(limit int) ParamsOptions {
	return func(params *ratelimiter.Params) {
		params.LimitByIP = true
		params.LimitIPBySecond = limit
	}
}

func WithLimitByToken(tokenInfo TokenInfo) ParamsOptions {
	return func(params *ratelimiter.Params) {
		params.LimitByToken = true
		for tokenID, limit := range tokenInfo {
			params.TokenList[tokenID] = limit
		}
	}
}

func WithBlockTimeToken(blockTimeToken time.Duration) ParamsOptions {
	return func(params *ratelimiter.Params) {
		params.BlockTimeToken = blockTimeToken
	}
}

func WithBlockTimeIP(blockTimeIP time.Duration) ParamsOptions {
	return func(params *ratelimiter.Params) {
		params.BlockTimeIP = blockTimeIP
	}
}

type Middleware struct {
	service ratelimiter.Service
}

func New(cfg DatabaseConfig, opts ...ParamsOptions) (*Middleware, error) {
	params := defaultParams()
	for _, opt := range opts {
		opt(&params)
	}

	cli := redis.NewClient(redis.Config{
		Addr:     cfg.Addr,
		Port:     cfg.Port,
		Password: cfg.Password,
	})
	blockedRepo := blocked.NewRepository(cli)
	accessRepo := access.NewRepository(cli)

	service, err := ratelimiter.NewService(accessRepo, blockedRepo, params)
	if err != nil {
		return nil, err
	}

	return &Middleware{
		service: service,
	}, nil
}

func (mw *Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("api_key")
		ip, err := httpInternal.GetIP(r)
		if err != nil {
			handleErr(w, ratelimiter.ErrGettingIP)
			return
		}

		err = mw.service.CanAccess(r.Context(), token, ip)
		if err != nil {
			handleErr(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func WithDatabaseConfig(addr string, port int, pwd string) DatabaseConfig {
	return DatabaseConfig{
		Addr:     addr,
		Port:     port,
		Password: pwd,
	}
}
