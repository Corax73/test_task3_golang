package customRedis

import (
	"checklist/customLog"
	"checklist/utils"
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	EnvPort, EnvPassword string
	RedisClient          *redis.Client
	Ctx                  context.Context
}

func GetClient(ctx context.Context) *RedisClient {
	envData := utils.GetConfFromEnvFile()
	redisPort := ":6379"
	redisPassword := ""
	if val, ok := envData["REDIS_PORT"]; ok {
		redisPort = utils.ConcatSlice([]string{":", val})
	}
	if val, ok := envData["REDIS_PASSWORD"]; ok {
		redisPassword = val
	}
	return &RedisClient{
		redisPort,
		redisPassword,
		redis.NewClient(&redis.Options{
			Addr:     utils.ConcatSlice([]string{"localhost", redisPort}),
			Password: redisPassword,
			DB:       0,
		}),
		ctx,
	}
}

func (redisClient *RedisClient) getAllKeys() []string {
	resp := make([]string, 0)
	result, err := redisClient.RedisClient.Keys(redisClient.Ctx, "*").Result()
	if err != nil {
		customLog.Logging(err)
	} else {
		resp = result
	}
	return resp
}

func (redisClient *RedisClient) RemoveModelKeys(modelName string) {
	if modelName != "" {
		keys := redisClient.getAllKeys()
		for _, val := range keys {
			if strings.HasPrefix(val, modelName) {
				redisClient.RedisClient.Del(redisClient.Ctx, val)
			}
		}
	}
}
