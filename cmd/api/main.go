package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/configuration"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {

	var cfg *configuration.Configuration
	profile := os.Getenv("profile")

	cfg, _ = configuration.LoadConfiguration(profile)
	fmt.Printf("%+v\n", cfg)

	uriMongoDB := "mongodb://" + cfg.DBConfiguration.MongoDBConfiguration.Username + ":" +
		cfg.DBConfiguration.MongoDBConfiguration.Password + "@" +
		cfg.DBConfiguration.MongoDBConfiguration.Host

	clientMongoDB, err := mongo.Connect(
		options.Client().ApplyURI(uriMongoDB),
	)

	if err != nil {
		panic(err)
	}

	databaseMongoDB := clientMongoDB.Database(cfg.DBConfiguration.MongoDBConfiguration.DBName)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.DBConfiguration.RedisConfiguration.Host,
		Password: cfg.DBConfiguration.RedisConfiguration.Password,
		DB:       cfg.DBConfiguration.RedisConfiguration.DBN,
	})

	imageKafkaProducer := NewImageProducerImpl(cfg.KafkaConfiguration.BootstrapServer)

	getImageUseCase := NewGetImageUseCase(databaseMongoDB, redisClient)
	insertImageUseCase := NewInsertImageUseCase(imageKafkaProducer, databaseMongoDB, redisClient)
	deleteImageUseCase := NewDeleteImageUseCase(databaseMongoDB)

	imageHandlers := NewImageHandler(
		getImageUseCase,
		insertImageUseCase,
		deleteImageUseCase,
	)

	chiRouter := chi.NewRouter()
	chiRouter.Use(DefaultHeaders)

	chiRouter.Route("/v1", func(r chi.Router) {
		for key, val := range imageHandlers.GetAllHandlers() {
			r.HandleFunc(key, val)
		}
	})

	err = http.ListenAndServe(cfg.SRVRConfiguration.Port, chiRouter)

	if err != nil {
		panic(err)
	}
}

func DefaultHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
