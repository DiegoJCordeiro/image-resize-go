package usecases

import (
	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/domain/models"
	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/domain/repositories"
	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/domain/storage"
)

type GetImageUseCase interface {
	Execute(uid string) (*models.ImageModel, error)
}

type GetImageUseCaseImpl struct {
	cacheRepository repositories.ImageCacheRepository
	noSqlRepository repositories.ImageNoSQLRepository
	s3Storage       storage.Storage
}

func NewGetImageUseCase(
	noSqlRepository repositories.ImageNoSQLRepository,
	cacheRepository repositories.ImageCacheRepository,
	s3Storage storage.Storage,
) GetImageUseCase {
	return &GetImageUseCaseImpl{
		cacheRepository: cacheRepository,
		noSqlRepository: noSqlRepository,
		s3Storage:       s3Storage,
	}
}

func (giuc *GetImageUseCaseImpl) Execute(uid string) (*models.ImageModel, error) {

	imageFound, err := giuc.getImageFromCache(uid)

	if err != nil {
		return nil, err
	}

	imageFound, err = giuc.getImageFromDatabase(uid)

	if err != nil {
		return nil, err
	}

	return imageFound, nil
}

func (giuc *GetImageUseCaseImpl) getImageFromCache(uid string) (*models.ImageModel, error) {

	imageFound, err := giuc.cacheRepository.GetImage(uid)

	if err != nil {
		return nil, err
	}

	if imageFound != nil {
		imageFound, _ = giuc.getObjectFromStorage(imageFound)
		return imageFound, nil
	}

	return imageFound, nil
}

func (giuc *GetImageUseCaseImpl) getImageFromDatabase(uid string) (*models.ImageModel, error) {

	imageFound, err := giuc.noSqlRepository.GetImage(uid)

	if err != nil {
		return nil, err
	}

	err = giuc.cacheRepository.PutImage(uid, imageFound)

	if err != nil {
		return nil, err
	}

	imageFound, _ = giuc.getObjectFromStorage(imageFound)

	return imageFound, nil
}

func (giuc *GetImageUseCaseImpl) getObjectFromStorage(imageFound *models.ImageModel) (*models.ImageModel, error) {

	imageObject, err := giuc.s3Storage.GetObject(imageFound.FullPath)

	if err != nil {
		return nil, err
	}

	imageFound.Object = imageObject

	return imageFound, nil
}
