package repository

import (
	"context"
	"line-translate-bot/pkg/crypto"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedisClient() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	opt, _ := redis.ParseURL(redisURL)
	return redis.NewClient(opt)
}

func NewRedisRepository(rdb *redis.Client) *RedisRepository {
	return &RedisRepository{
		rdb: rdb,
		ctx: context.Background(),
	}
}

func (r *RedisRepository) SetGroupLanguage(groupID string, lang string) error {
	key := []byte(os.Getenv("ENCRYPTION_KEY"))
	cipherText, err := crypto.Encrypt([]byte(lang), key)
	if err != nil {
		return err
	}
	return r.rdb.Set(r.ctx, groupID+":lang", cipherText, 0).Err()
}

func (r *RedisRepository) GetGroupLanguage(groupID string) (string, error) {
	key := []byte(os.Getenv("ENCRYPTION_KEY"))
	encryptedData, err := r.rdb.Get(r.ctx, groupID+":lang").Bytes()
	if err != nil {
		return "", err
	}
	plainText, err := crypto.Decrypt(encryptedData, key)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

func (r *RedisRepository) SaveLanguages(sourceID string, langs []string) error {
	key := []byte(os.Getenv("ENCRYPTION_KEY"))

	// 將 []string 轉為一個單一字串 (例如用逗號分隔)
	// 這裡我們直接用 fmt.Sprint 或是 strings.Join，根據你的需求
	langsStr := strings.Join(langs, ",")

	cipherText, err := crypto.Encrypt([]byte(langsStr), key)
	if err != nil {
		return err
	}

	return r.rdb.Set(r.ctx, sourceID+":langs", cipherText, 0).Err()
}

// GetLanguages 讀取該群組支援的語言清單
func (r *RedisRepository) GetLanguages(sourceID string) ([]string, error) {
	key := []byte(os.Getenv("ENCRYPTION_KEY"))

	encryptedData, err := r.rdb.Get(r.ctx, sourceID+":langs").Bytes()
	if err == redis.Nil {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}

	plainText, err := crypto.Decrypt(encryptedData, key)
	if err != nil {
		return nil, err
	}

	langsStr := string(plainText)
	return strings.Split(langsStr, ","), nil
}
