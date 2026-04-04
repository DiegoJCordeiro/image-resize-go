//go:build wireinject
// +build wireinject

package main

import (
	"github.com/DiegoJCordeiro/image-resizing-go-producer-api/internal/domain/events"
	"github.com/DiegoJCordeiro/image-resizing-go-producer-api/internal/domain/repositories"
	"github.com/DiegoJCordeiro/image-resizing-go-producer-api/internal/domain/storage"
	"github.com/DiegoJCordeiro/image-resizing-go-producer-api/internal/infra/cache"
	"github.com/DiegoJCordeiro/image-resizing-go-producer-api/internal/infra/http/handlers"
	"github.com/DiegoJCordeiro/image-resizing-go-producer-api/internal/infra/kafka"
	"github.com/DiegoJCordeiro/image-resizing-go-producer-api/internal/infra/persistence"
	"github.com/DiegoJCordeiro/image-resizing-go-producer-api/internal/infra/s3storage"
	"github.com/DiegoJCordeiro/image-resizing-go-producer-api/internal/usecases"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	setKafkaProducerDependency = wire.NewSet(
		kafka.NewImageProducerImpl,
		wire.Bind(new(events.EventProducer), new(*kafka.ImageKafkaProducer)),
	)

	setRedisRepositoryDependency = wire.NewSet(
		cache.NewImageRedisBRepositoryImpl,
		wire.Bind(new(repositories.ImageCacheRepository), new(*cache.ImageRedisBRepositoryImpl)),
	)

	setMongoDBRepositoryDependency = wire.NewSet(
		persistence.NewImageMongoDBRepositoryImpl,
		wire.Bind(new(repositories.ImageNoSQLRepository), new(*persistence.ImageMongoDBRepositoryImpl)),
	)

	setS3Dependency = wire.NewSet(
		s3storage.NewS3StorageImpl,
		wire.Bind(new(storage.Storage), new(*s3storage.S3StorageImpl)),
	)
)

// Handlers
func NewImageHandler(
	getImageUseCase usecases.GetImageUseCase,
	insertImageUseCase usecases.InsertImageUseCase,
	deleteImageUseCase usecases.DeleteImageUseCase,
) *handlers.ImageHandler {

	wire.Build(
		handlers.NewImageHandler,
	)

	return &handlers.ImageHandler{}
}

// UseCases
func NewGetImageUseCase(database *mongo.Database, redisClient *redis.Client) usecases.GetImageUseCase {

	wire.Build(
		setS3Dependency,
		setRedisRepositoryDependency,
		setMongoDBRepositoryDependency,
		usecases.NewGetImageUseCase,
	)

	return &usecases.GetImageUseCaseImpl{}
}

func NewInsertImageUseCase(kafkaProducer events.EventProducer, database *mongo.Database, redisClient *redis.Client) usecases.InsertImageUseCase {

	wire.Build(
		setS3Dependency,
		setRedisRepositoryDependency,
		setMongoDBRepositoryDependency,
		usecases.NewInsertImageUseCase,
	)

	return &usecases.InsertImageUseCaseImpl{}
}

func NewDeleteImageUseCase(database *mongo.Database, redisClient *redis.Client) usecases.DeleteImageUseCase {

	wire.Build(
		setMongoDBRepositoryDependency,
		setRedisRepositoryDependency,
		usecases.NewDeleteImageUseCase,
	)

	return &usecases.DeleteImageUseCaseImpl{}
}

// Producers
func NewImageProducerImpl(bootstrapServer string) *kafka.ImageKafkaProducer {

	wire.Build(
		kafka.NewImageProducerImpl,
	)

	return &kafka.ImageKafkaProducer{}
}
