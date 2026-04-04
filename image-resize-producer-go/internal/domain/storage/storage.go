package storage

import "github.com/DiegoJCordeiro/image-resizing-go-producer-api/internal/domain/models"

type Storage interface {
	GetObject(fullPath string) ([]byte, error)
	PutObject(imageModel *models.ImageModel) error
}
