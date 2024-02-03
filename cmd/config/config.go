package config

import (
	"encoding/json"

	"github.com/mateusmatinato/goexpert-rate-limiter/internal/platform/redis"
	"github.com/spf13/viper"
)

type Config struct {
	RedisURL      string `mapstructure:"REDIS_URL"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisPort     int    `mapstructure:"REDIS_PORT"`
	TokenInfo     map[string]TokenInfo
}

type TokenInfo struct {
	TokenID            string `json:"token_id"`
	MaxRequestsSeconds int    `json:"max_requests_seconds"`
	BlockTimeSeconds   int    `json:"block_time_seconds"`
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
	cfg.TokenInfo = make(map[string]TokenInfo)
	for _, token := range tokenList {
		cfg.TokenInfo[token.TokenID] = token
	}

	return
}
