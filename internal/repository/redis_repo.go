package repository

import (
	"context"
	"encoding/json"
	"line-translate-bot/pkg/crypto"
	"os"

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
	// 將陣列轉為 JSON 字串
	data, err := json.Marshal(langs)
	if err != nil {
		return err
	}
	// 加密 JSON 字串
	key := []byte(os.Getenv("ENCRYPTION_KEY"))
	cipherText, err := crypto.Encrypt(data, key)
	if err != nil {
		return err
	}
	return r.rdb.Set(r.ctx, sourceID+":langs", cipherText, 0).Err()
}

// GetLanguages 讀取該群組支援的語言清單
func (r *RedisRepository) GetLanguages(sourceID string) ([]string, error) {
	encryptedData, err := r.rdb.Get(r.ctx, sourceID+":langs").Bytes()
	if err != nil {
		return nil, err
	}

	key := []byte(os.Getenv("ENCRYPTION_KEY"))
	plainText, err := crypto.Decrypt(encryptedData, key)
	if err != nil {
		return nil, err
	}

	// 從 JSON 還原回 []string
	var langs []string
	err = json.Unmarshal(plainText, &langs)
	return langs, err
}
