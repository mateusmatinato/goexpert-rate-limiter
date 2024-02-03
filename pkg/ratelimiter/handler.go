package ratelimiter

import (
	"errors"
	"net/http"

	"github.com/mateusmatinato/goexpert-rate-limiter/internal/ratelimiter"
)

func handleErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ratelimiter.ErrInvalidToken):
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
	case errors.Is(err, ratelimiter.ErrBlocked):
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(err.Error()))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
