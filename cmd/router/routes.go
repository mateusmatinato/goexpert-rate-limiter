package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mateusmatinato/goexpert-rate-limiter/cmd/config"
	"github.com/mateusmatinato/goexpert-rate-limiter/pkg/ratelimiter"
)

func StartTestRoutes(cfg config.Config) *mux.Router {
	limiterByToken, err := ratelimiter.New(
		ratelimiter.WithDatabaseConfig(cfg.RedisURL, cfg.RedisPort, cfg.RedisPassword),
		ratelimiter.WithLimitByToken(newTokenInfo(cfg)),
		ratelimiter.WithBlockTimeToken(cfg.BlockTimeToken),
	)
	if err != nil {
		panic("error starting rate limiter by token: " + err.Error())
	}

	limiterByIP, err := ratelimiter.New(
		ratelimiter.WithDatabaseConfig(cfg.RedisURL, cfg.RedisPort, cfg.RedisPassword),
		ratelimiter.WithLimitByIP(cfg.LimitByIP),
		ratelimiter.WithBlockTimeIP(cfg.BlockTimeIP),
	)
	if err != nil {
		panic("error starting rate limiter by ip: " + err.Error())
	}

	limiterByBoth, err := ratelimiter.New(
		ratelimiter.WithDatabaseConfig(cfg.RedisURL, cfg.RedisPort, cfg.RedisPassword),
		ratelimiter.WithLimitByToken(newTokenInfo(cfg)),
		ratelimiter.WithBlockTimeToken(cfg.BlockTimeToken),
		ratelimiter.WithLimitByIP(cfg.LimitByIP),
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
	r.Handle("/", limiterByBoth.Middleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Testing Limiter Default Route"))
		}),
	))
	return r
}

func newTokenInfo(cfg config.Config) ratelimiter.TokenInfo {
	tokenInfo := make(ratelimiter.TokenInfo)
	for _, token := range cfg.TokenList {
		tokenInfo[token.ID] = token.RequestLimitSecond
	}
	return tokenInfo
}
