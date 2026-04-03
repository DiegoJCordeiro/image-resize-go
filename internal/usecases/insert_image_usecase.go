package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/domain/events"
	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/domain/models"
	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/domain/models/status"
	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/domain/repositories"
	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/domain/storage"
	"github.com/google/uuid"
)

var (
	extensionsAllowed = [4]string{
		".jpg", ".png", ".jpeg", ".gif",
	}
)

type InsertImageUseCase interface {
	Execute(filename, extension string, processType string, object []byte) (string, error)
}

type InsertImageUseCaseImpl struct {
	kafkaProducer   events.EventProducer
	cacheRepository repositories.ImageCacheRepository
	noSqlRepository repositories.ImageNoSQLRepository
	s3Storage       storage.Storage
}

func NewInsertImageUseCase(
	kafkaProducer events.EventProducer,
	noSqlRepository repositories.ImageNoSQLRepository,
	cacheRepository repositories.ImageCacheRepository,
	s3Storage storage.Storage,
) InsertImageUseCase {
	return &InsertImageUseCaseImpl{
		cacheRepository: cacheRepository,
		noSqlRepository: noSqlRepository,
		kafkaProducer:   kafkaProducer,
		s3Storage:       s3Storage,
	}
}

func (iuc *InsertImageUseCaseImpl) Execute(filename, extension string, processType string, object []byte) (string, error) {

	if iuc.isExtensionValid(extension) == false {
		return "", fmt.Errorf("extension %s is invalid", extension)
	}

	ctxParent := context.Background()
	ctxTimeout, cancelFunc := context.WithTimeout(ctxParent, time.Second*2)

	defer cancelFunc()

	errChan := make(chan error, 1)
	imageUUID := uuid.New()

	prefixS3 := fmt.Sprintf(
		"%s/original/%s/%s", imageUUID.String(), processType, filename)

	imageModel := models.NewImageModel(
		imageUUID.String(),
		filename,
		extension,
		prefixS3,
		object,
		models.ProcessType(processType),
		status.InProgress,
	)

	go iuc.publishAtStorage(imageModel, errChan)
	go iuc.persistenceAtDatabase(imageModel, errChan)
	go iuc.persistenceAtCache(imageModel, errChan)
	go iuc.publishEventAtKafka(imageModel, errChan)

	select {
	case err := <-errChan:
		return "", fmt.Errorf("[InsertImage] - [Error]: %s", err.Error())
	case <-ctxTimeout.Done():
		return imageUUID.String(), nil
	}
}

func (iuc *InsertImageUseCaseImpl) isExtensionValid(extension string) bool {

	for _, allowed := range extensionsAllowed {
		if extension == allowed {
			return true
		}
	}

	return false
}

func (iuc *InsertImageUseCaseImpl) persistenceAtCache(imageModel *models.ImageModel, chanErr chan<- error) {

	err := iuc.cacheRepository.PutImage(imageModel.UID, imageModel)

	if err != nil {
		chanErr <- err
	}
}

func (iuc *InsertImageUseCaseImpl) persistenceAtDatabase(imageModel *models.ImageModel, chanErr chan<- error) {

	err := iuc.noSqlRepository.InsertImage(imageModel)

	if err != nil {
		chanErr <- err
	}
}

func (iuc *InsertImageUseCaseImpl) publishEventAtKafka(imageModel *models.ImageModel, chanErr chan<- error) {

	jsonPayload, err := json.Marshal(imageModel)

	if err != nil {
		chanErr <- err
	}

	err = iuc.kafkaProducer.Publish(string(jsonPayload), "image-process", []byte(imageModel.ProcessType))

	if err != nil {
		chanErr <- err
	}
}

func (iuc *InsertImageUseCaseImpl) publishAtStorage(imageModel *models.ImageModel, chanErr chan<- error) {

	err := iuc.s3Storage.PutObject(imageModel)

	if err != nil {
		chanErr <- err
	}
}
