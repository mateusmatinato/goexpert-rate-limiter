package config

import (
	"encoding/json"
	"time"

	"github.com/mateusmatinato/goexpert-rate-limiter/internal/platform/redis"
	"github.com/spf13/viper"
)

type Config struct {
	RedisURL       string        `mapstructure:"REDIS_URL"`
	RedisPassword  string        `mapstructure:"REDIS_PASSWORD"`
	RedisPort      int           `mapstructure:"REDIS_PORT"`
	LimitByIP      int           `mapstructure:"LIMIT_BY_IP"`
	BlockTimeIP    time.Duration `mapstructure:"BLOCK_TIME_IP"`
	BlockTimeToken time.Duration `mapstructure:"BLOCK_TIME_TOKEN"`
	TokenList      []TokenInfo
}

type TokenInfo struct {
	ID                 string `json:"id"`
	RequestLimitSecond int    `json:"request_limit_second"`
}

func (c *Config) ToRedisConfig() redis.Config {
	return redis.Config{
		Addr:     c.RedisURL,
		Port:     c.RedisPort,
		Password: c.RedisPassword,
	}
}

func LoadConfig(path string) (cfg Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return
	}

	tokenCfg := viper.GetString("TOKEN_INFO")
	var tokenList []TokenInfo
	err = json.Unmarshal([]byte(tokenCfg), &tokenList)
	if err != nil {
		return
	}
	cfg.TokenList = tokenList

	return
}
