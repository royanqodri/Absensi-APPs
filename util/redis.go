package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/royanqodri/Absensi/util/constants"
)

var RedisClient *redis.Client
var bgCtx = context.Background() // nama lebih jelas

func InitRedis() {
	// cfg := config.Get().Redis
	// addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	// RedisClient = redis.NewClient(&redis.Options{
	// 	Addr:     addr,
	// 	Password: cfg.Password,
	// 	DB:       cfg.DB,
	// })

	_, err := RedisClient.Ping(bgCtx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

// ──────────────────────────────────────────────────────────────
// Helper: ambil context yang benar (fallback ke background)
func getRedisCtx(c echo.Context) context.Context {
	if c == nil {
		return bgCtx
	}
	return c.Request().Context() // ← yang benar!
}

// ──────────────────────────────────────────────────────────────

func IsKeyExistInRedis(c echo.Context, customerNo, tableName string) (bool, error) {
	key := strings.ToLower(fmt.Sprintf("%s:%s:%s", constants.REDIS_KEY_MASTER_DATA, customerNo, tableName))

	ctx := getRedisCtx(c)
	count, err := RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func StoreHashRecordRedis(c echo.Context, customerNo, tableName, fieldKey string, data interface{}) error {
	key := strings.ToLower(fmt.Sprintf("%s:%s:%s", constants.REDIS_KEY_MASTER_DATA, customerNo, tableName))

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ctx := getRedisCtx(c)
	_, err = RedisClient.HSet(ctx, key, strings.ToLower(fieldKey), string(b)).Result()
	return err
}

func GetHashRecordRedis(c echo.Context, customerNo, tableName, fieldKey string) (string, error) {
	key := strings.ToLower(fmt.Sprintf("%s:%s:%s", constants.REDIS_KEY_MASTER_DATA, customerNo, tableName))

	ctx := getRedisCtx(c)
	val, err := RedisClient.HGet(ctx, key, strings.ToLower(fieldKey)).Result()
	return val, err
}

func GetHashAllRecordsRedis(c echo.Context, customerNo, tableName string) (map[string]string, error) {
	key := strings.ToLower(fmt.Sprintf("%s:%s:%s", constants.REDIS_KEY_MASTER_DATA, customerNo, tableName))

	ctx := getRedisCtx(c)
	return RedisClient.HGetAll(ctx, key).Result()
}

func DeleteHashRecordRedis(c echo.Context, customerNo, tableName, fieldKey string) error {
	key := strings.ToLower(fmt.Sprintf("%s:%s:%s", constants.REDIS_KEY_MASTER_DATA, customerNo, tableName))

	ctx := getRedisCtx(c)
	_, err := RedisClient.HDel(ctx, key, strings.ToLower(fieldKey)).Result()
	return err
}

func DeleteKeyRedis(c echo.Context, customerNo, tableName string) error {
	key := strings.ToLower(fmt.Sprintf("%s:%s:%s", constants.REDIS_KEY_MASTER_DATA, customerNo, tableName))

	ctx := getRedisCtx(c)
	_, err := RedisClient.Del(ctx, key).Result()
	return err
}
