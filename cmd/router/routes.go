package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mateusmatinato/goexpert-rate-limiter/cmd/config"
	"github.com/mateusmatinato/goexpert-rate-limiter/pkg/ratelimiter"
)

func StartTestRoutes(cfg config.Config) *mux.Router {
	limiterByToken, err := ratelimiter.NewRateLimiterMiddleware(
		newRateLimiterDBConfig(cfg),
		ratelimiter.WithBlockByToken(newTokenList(cfg)),
		ratelimiter.WithBlockTimeToken(cfg.BlockTimeToken),
	)
	if err != nil {
		panic("error starting rate limiter by token: " + err.Error())
	}

	limiterByIP, err := ratelimiter.NewRateLimiterMiddleware(
		newRateLimiterDBConfig(cfg),
		ratelimiter.WithBlockByIP(cfg.LimitByIP),
		ratelimiter.WithBlockTimeIP(cfg.BlockTimeIP),
	)
	if err != nil {
		panic("error starting rate limiter by ip: " + err.Error())
	}

	limiterByBoth, err := ratelimiter.NewRateLimiterMiddleware(
		newRateLimiterDBConfig(cfg),
		ratelimiter.WithBlockByToken(newTokenList(cfg)),
		ratelimiter.WithBlockTimeToken(cfg.BlockTimeToken),
		ratelimiter.WithBlockByIP(cfg.LimitByIP),
		ratelimiter.WithBlockTimeIP(cfg.BlockTimeIP),
	)
	if err != nil {
		panic("error starting rate limiter by both: " + err.Error())
	}

	r := mux.NewRouter()
	r.Handle("/token", limiterByToken.Middleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Testing Limiter By TOKEN"))
		}),
	))
	r.Handle("/ip", limiterByIP.Middleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Testing Limiter By IP"))
		}),
	))
	r.Handle("/both", limiterByBoth.Middleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Testing Limiter By BOTH"))
		}),
	))
	r.Handle("", limiterByBoth.Middleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Testing Limiter Default Route"))
		}),
	))
	return r
}

func newTokenList(cfg config.Config) []ratelimiter.TokenInfo {
	var tokenList []ratelimiter.TokenInfo
	for _, token := range cfg.TokenList {
		tokenList = append(tokenList, ratelimiter.TokenInfo{
			ID:                  token.ID,
			MaxRequestsBySecond: token.RequestLimitSecond,
		})
	}
	return tokenList
}

func newRateLimiterDBConfig(cfg config.Config) ratelimiter.DatabaseConfig {
	return ratelimiter.DatabaseConfig{
		Addr:     cfg.RedisURL,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
	}
}
