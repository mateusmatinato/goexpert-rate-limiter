package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/mateusmatinato/goexpert-rate-limiter/cmd/config"
	"github.com/mateusmatinato/goexpert-rate-limiter/cmd/router"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	r *mux.Router
}

func (s *TestSuite) SetupTest() {
	cfg := config.Config{
		RedisURL:       "localhost",
		RedisPassword:  "user",
		RedisPort:      6379,
		BlockTimeToken: 5 * time.Second,
		BlockTimeIP:    5 * time.Second,
		TokenList: []config.TokenInfo{
			{
				ID:                 "test",
				RequestLimitSecond: 5,
			},
			{
				ID:                 "test2",
				RequestLimitSecond: 5,
			},
		},
		LimitByIP: 5,
	}
	r := router.StartTestRoutes(cfg)

	s.r = r
}

func TestRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestTokenLimit() {
	s.Run("should block after 5 requests and unblock after 5 seconds", func() {
		req, _ := http.NewRequest("GET", "/token", nil)
		req.Header.Set("api_key", "test")
		req.Header.Set("X-Forwarded-For", "127.0.0.1")

		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			s.r.ServeHTTP(w, req)
			s.Equal(http.StatusOK, w.Code)
		}

		w := httptest.NewRecorder()
		s.r.ServeHTTP(w, req)
		s.Equal(http.StatusTooManyRequests, w.Code)

		time.Sleep(6 * time.Second)

		w = httptest.NewRecorder()
		s.r.ServeHTTP(w, req)
		s.Equal(http.StatusOK, w.Code)
	})
}

func (s *TestSuite) TestIPLimit() {
	s.Run("should block after 5 requests and unblock after 5 seconds", func() {
		req, _ := http.NewRequest("GET", "/ip", nil)
		req.Header.Set("X-Forwarded-For", "127.0.0.2")

		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			s.r.ServeHTTP(w, req)
			s.Equal(http.StatusOK, w.Code)
		}

		w := httptest.NewRecorder()
		s.r.ServeHTTP(w, req)
		s.Equal(http.StatusTooManyRequests, w.Code)

		time.Sleep(6 * time.Second)

		w = httptest.NewRecorder()
		s.r.ServeHTTP(w, req)
		s.Equal(http.StatusOK, w.Code)
	})
}

func (s *TestSuite) TestTokenAndIPLimit() {
	s.Run("should block by ip when api_key is empty and allow if token is sent after", func() {
		req, _ := http.NewRequest("GET", "/both", nil)
		req.Header.Set("api_key", "")
		req.Header.Set("X-Forwarded-For", "127.0.0.3")

		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			s.r.ServeHTTP(w, req)
			s.Equal(http.StatusOK, w.Code)
		}

		w := httptest.NewRecorder()
		s.r.ServeHTTP(w, req)
		s.Equal(http.StatusTooManyRequests, w.Code)

		req.Header.Set("api_key", "test2")
		w = httptest.NewRecorder()
		s.r.ServeHTTP(w, req)
		s.Equal(http.StatusOK, w.Code)
	})
}
