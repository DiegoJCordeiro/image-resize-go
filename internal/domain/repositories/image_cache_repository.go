package repositories

import "github.com/DiegoJCordeiro/image-resizing-go-api/internal/domain/models"

type ImageCacheRepository interface {
	GetImage(key string) (*models.ImageModel, error)
	PutImage(key string, value *models.ImageModel) error
	DeleteImage(key string) error
}
