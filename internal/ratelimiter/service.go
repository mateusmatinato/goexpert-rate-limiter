package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/mateusmatinato/goexpert-rate-limiter/internal/platform/log"
	"github.com/mateusmatinato/goexpert-rate-limiter/internal/ratelimiter/access"
	"github.com/mateusmatinato/goexpert-rate-limiter/internal/ratelimiter/blocked"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks
type Service interface {
	CanAccess(ctx context.Context, token string, ip string) error
}

type service struct {
	accessRepository  access.Repository
	blockedRepository blocked.Repository
	params            Params
}

func (s service) CanAccess(ctx context.Context, token string, ip string) error {
	if s.params.BlockByToken && token != "" {
		limit, ok := s.params.TokenList[token]
		if !ok {
			return ErrInvalidToken
		}

		return s.validateAccess(ctx, OriginToken, token, limit)
	}

	if s.params.BlockByIP {
		return s.validateAccess(ctx, OriginIP, ip, s.params.LimitIPBySecond)
	}

	if token == "" {
		return ErrInvalidSetup
	}

	return nil
}

func (s service) validateAccess(ctx context.Context, origin LimitOrigin, key string, limit int) error {
	isBlocked, err := s.blockedRepository.IsBlocked(ctx, generateKey(origin, key))
	if err != nil {
		log.Error("error validating is blocked", err)
		return ErrInternal
	}

	if isBlocked {
		log.Info("access is blocked", fmt.Sprintf("origin:%s", origin), fmt.Sprintf("key:%s", key))
		return ErrBlocked
	}

	count, err := s.accessRepository.GetAccessCount(ctx, generateKey(origin, key))
	if err != nil {
		log.Error("error getting access count", err)
		return ErrInternal
	}

	if count >= limit {
		log.Info("blocking access", fmt.Sprintf("origin:%s", origin), fmt.Sprintf("key:%s", key))
		err := s.blockedRepository.Block(ctx, generateKey(origin, key), getBlockDuration(s.params, origin))
		if err != nil {
			log.Error("error blocking access", err)
			return ErrInternal
		}
		return ErrBlocked
	}

	err = s.accessRepository.IncrementAccessCount(ctx, generateKey(origin, key))
	if err != nil {
		log.Error("error incrementing access", err)
		return ErrInternal
	}
	log.Info("access allowed", fmt.Sprintf("origin:%s", origin),
		fmt.Sprintf("key:%s", key), fmt.Sprintf("count:%d", count+1))
	return nil
}

func getBlockDuration(params Params, origin LimitOrigin) time.Duration {
	if origin == OriginToken {
		return params.BlockTimeToken
	}

	return params.BlockTimeIP
}

func NewService(accessRepository access.Repository, blockedRepository blocked.Repository, params Params) (Service, error) {
	err := validateParams(params)
	if err != nil {
		return nil, err
	}

	return &service{
		accessRepository:  accessRepository,
		blockedRepository: blockedRepository,
		params:            params,
	}, nil
}

func validateParams(params Params) error {
	if !params.BlockByIP && !params.BlockByToken {
		return ErrNotBlocking
	}

	if params.BlockByToken && len(params.TokenList) == 0 {
		return ErrTokenListEmpty
	}

	if params.BlockByIP && params.LimitIPBySecond <= 0 {
		return ErrInvalidLimitIP
	}

	if params.BlockByToken {
		for _, limit := range params.TokenList {
			if limit <= 0 {
				return ErrInvalidLimitToken
			}
		}
	}

	return nil
}
