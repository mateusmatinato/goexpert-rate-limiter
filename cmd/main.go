package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mateusmatinato/goexpert-rate-limiter/cmd/config"
	"github.com/mateusmatinato/goexpert-rate-limiter/pkg/ratelimiter"
)

func main() {
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		panic(fmt.Sprintf("error starting configs: %s", err.Error()))
	}

	// TODO: read from configs and set
	limiterMW, err := ratelimiter.NewRateLimiterMiddleware(ratelimiter.DatabaseConfig{
		Addr:     cfg.RedisURL,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
	}, ratelimiter.WithBlockByToken([]ratelimiter.TokenInfo{
		{
			ID:                  "test",
			MaxRequestsByMinute: 5,
		},
		{
			ID:                  "test2",
			MaxRequestsByMinute: 10,
		},
	},
	), ratelimiter.WithBlockTimeToken(10*time.Second))
	if err != nil {
		panic(fmt.Sprintf("error starting rate limiter: %s", err.Error()))
	}

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	r.Use(limiterMW.Middleware)

	http.ListenAndServe(":8080", r)
}
