package redis

import (
	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/config"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func NewRedis(config *config.Config) *redis.Redis {
	return redis.MustNewRedis(config.RedisConf)
}
