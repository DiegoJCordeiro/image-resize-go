package usecases

import (
	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/domain/models/status"
	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/domain/repositories"
)

type DeleteImageUseCase interface {
	Execute(uid string) error
}

type DeleteImageUseCaseImpl struct {
	cacheRepository repositories.ImageCacheRepository
	noSQLRepository repositories.ImageNoSQLRepository
}

func NewDeleteImageUseCase(
	cacheRepository repositories.ImageCacheRepository,
	noSQLRepository repositories.ImageNoSQLRepository,
) DeleteImageUseCase {
	return &DeleteImageUseCaseImpl{
		cacheRepository: cacheRepository,
		noSQLRepository: noSQLRepository,
	}
}

func (diuc *DeleteImageUseCaseImpl) Execute(uid string) error {

	err := diuc.noSQLRepository.DeleteImage(uid, string(status.Deleted))

	if err != nil {
		return err
	}

	err = diuc.cacheRepository.DeleteImage(uid)

	if err != nil {
		return err
	}

	return nil
}
