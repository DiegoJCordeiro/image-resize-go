package persistence

import (
	"context"
	"time"

	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/domain/models"
	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/infra/persistence/collections"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	collectionName = "image_resizing"
)

type ImageMongoDBRepositoryImpl struct {
	mongoCollection *mongo.Collection
}

// NewImageMongoDBRepositoryImpl creates a new repository for mongodb
func NewImageMongoDBRepositoryImpl(database *mongo.Database) *ImageMongoDBRepositoryImpl {

	mongoCollection := database.Collection(collectionName)

	return &ImageMongoDBRepositoryImpl{
		mongoCollection: mongoCollection,
	}
}

// GetImage getting a image by uid
func (imr *ImageMongoDBRepositoryImpl) GetImage(uid string) (*models.ImageModel, error) {

	var imageCollection collections.ImageCollection

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelFunc()

	filter := bson.D{{Key: "image_id", Value: uid}}

	if err := imr.mongoCollection.FindOne(ctx, filter).Decode(&imageCollection); err != nil {
		return nil, err
	}

	return imageCollection.ToModel(), nil
}

// GetImageByUIDAndIsNotStatus getting an image by uid and status
func (imr *ImageMongoDBRepositoryImpl) GetImageByUIDAndIsNotStatus(uid string, status string) (*models.ImageModel, error) {

	var imageCollection collections.ImageCollection

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelFunc()

	filter := bson.D{
		{Key: "image_id", Value: uid},
		{Key: "status", Value: bson.D{
			{Key: "$ne", Value: status},
		}},
	}

	if err := imr.mongoCollection.FindOne(ctx, filter).Decode(&imageCollection); err != nil {
		return nil, err
	}

	return imageCollection.ToModel(), nil
}

// InsertImage deleting a image by uid
func (imr *ImageMongoDBRepositoryImpl) InsertImage(value *models.ImageModel) error {

	uid := bson.NewObjectID()
	createdAt := time.Now()

	imageCollection := collections.NewImageCollection(
		uid,
		value.UID,
		value.FileName,
		value.Extension,
		value.FullPath,
		string(value.ProcessType),
		string(value.Status),
		&createdAt,
		nil,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	_, err := imr.mongoCollection.InsertOne(ctx, imageCollection)

	if err != nil {
		panic(err)
	}

	return err
}

// DeleteImage deleting a image by uid
func (imr *ImageMongoDBRepositoryImpl) DeleteImage(uid string, status string) error {

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelFunc()

	filter := bson.D{{Key: "image_id", Value: uid}}
	fieldsToUpdate := bson.D{{Key: "$set", Value: bson.D{
		{Key: "status", Value: status},
		{Key: "deleted_at", Value: time.Now()},
	}}}

	_, err := imr.mongoCollection.UpdateOne(ctx, filter, fieldsToUpdate)

	if err != nil {
		return err
	}

	return nil
}
