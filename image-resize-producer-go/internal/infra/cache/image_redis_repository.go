package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/DiegoJCordeiro/image-resizing-go-producer-api/internal/domain/models"
	"github.com/redis/go-redis/v9"
)

const (
	redisPrefix = "image_resized:"
)

type ImageRedisBRepositoryImpl struct {
	redisClient *redis.Client
	ttl         time.Duration
}

func NewImageRedisBRepositoryImpl(redisClient *redis.Client) *ImageRedisBRepositoryImpl {
	return &ImageRedisBRepositoryImpl{
		redisClient: redisClient,
	}
}

func (imr *ImageRedisBRepositoryImpl) GetImage(key string) (*models.ImageModel, error) {

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	val, err := imr.redisClient.Get(ctx, redisPrefix+key).Result()

	if err != nil {

		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, err
	}

	var imageFound models.ImageModel

	if err := json.Unmarshal([]byte(val), &imageFound); err != nil {
		return nil, err
	}

	return &imageFound, nil
}

func (imr *ImageRedisBRepositoryImpl) PutImage(key string, value *models.ImageModel) error {

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	valueMarshaled, err := json.Marshal(value)

	if err != nil {
		return err
	}

	imr.redisClient.Set(ctx, redisPrefix+key, valueMarshaled, imr.ttl)

	return nil
}

func (imr *ImageRedisBRepositoryImpl) DeleteImage(key string) error {

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	imr.redisClient.Del(ctx, redisPrefix+key)

	return nil
}
