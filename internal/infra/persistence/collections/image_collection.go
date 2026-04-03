package collections

import (
	"time"

	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/domain/models"
	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/domain/models/status"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ImageCollection struct {
	UID          bson.ObjectID `bson:"_id, omitempty"`
	ImageID      string        `bson:"image_id, omitempty"`
	FileName     string        `bson:"filename, omitempty"`
	Extension    string        `bson:"extension, omitempty"`
	RelativePath string        `bson:"relative_path, omitempty"`
	Status       string        `bson:"status, omitempty"`
	ProcessType  string        `bson:"process_type, omitempty"`
	CreatedAt    *time.Time    `bson:"created_at, omitempty"`
	DeletedAt    *time.Time    `bson:"deleted_at, omitempty"`
}

func NewImageCollection(uid bson.ObjectID, imageID, fileName, extension, relativePath, processType, status string, createdAt, deletedAt *time.Time) *ImageCollection {
	return &ImageCollection{
		UID:          uid,
		ImageID:      imageID,
		FileName:     fileName,
		Extension:    extension,
		RelativePath: relativePath,
		Status:       status,
		ProcessType:  processType,
		CreatedAt:    createdAt,
		DeletedAt:    deletedAt,
	}
}

func (ic *ImageCollection) ToModel() *models.ImageModel {

	return &models.ImageModel{
		UID:         ic.ImageID,
		FileName:    ic.FileName,
		Extension:   ic.Extension,
		FullPath:    ic.RelativePath,
		Status:      status.Status(ic.Status),
		ProcessType: models.ProcessType(ic.ProcessType),
	}
}
