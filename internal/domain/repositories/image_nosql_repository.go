package repositories

import "github.com/DiegoJCordeiro/image-resizing-go-api/internal/domain/models"

type ImageNoSQLRepository interface {
	GetImage(uid string) (*models.ImageModel, error)
	InsertImage(value *models.ImageModel) error
	DeleteImage(uid string) error
}
