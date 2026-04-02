package usecases

import "github.com/DiegoJCordeiro/image-resizing-go-api/internal/domain/repositories"

type DeleteImageUseCase interface {
	Execute(uid string) error
}

type DeleteImageUseCaseImpl struct {
	cacheRepository repositories.ImageCacheRepository
	noSQLRepository repositories.ImageNoSQLRepository
}

func NewDeleteImageUseCase(
	noSQLRepository repositories.ImageNoSQLRepository,
) DeleteImageUseCase {
	return &DeleteImageUseCaseImpl{
		noSQLRepository: noSQLRepository,
	}
}

func (diuc *DeleteImageUseCaseImpl) Execute(uid string) error {

	err := diuc.noSQLRepository.DeleteImage(uid)

	if err != nil {
		return err
	}

	return nil
}
