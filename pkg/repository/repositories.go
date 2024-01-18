package repository

import (
	"github.com/go-redis/redis"
)

type Authorization interface {
}

type RedisRepository struct {
	client *redis.Client
}
type Repository struct {
	Authorization
}

func NewRepository() *Repository {
	return &Repository{}
}

func NewRedisRepository(host, port, password string) *RedisRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0, // use default DB
	})

	return &RedisRepository{
		client: rdb,
	}
}

// Метод для получения данных из Redis
func (r *RedisRepository) GetData(key string) ([]string, error) {
	zRange := r.client.ZRange(key, 0, -1)
	if zRange.Err() != nil {
		return nil, zRange.Err()
	}

	var data []string
	for _, v := range zRange.Val() {
		data = append(data, v)
	}

	return data, nil
}
