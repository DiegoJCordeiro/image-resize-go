package handlersdto

type ImageInsertedDTO struct {
	ImageUid string `json:"image-identifier"`
}

func NewImageInsertedDTO(imageUid string) *ImageInsertedDTO {
	return &ImageInsertedDTO{
		ImageUid: imageUid,
	}
}
