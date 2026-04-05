package handlersdto

type ImageInsertedDTO struct {
	ImageUid string `json:"iuid"`
}

func NewImageInsertedDTO(imageUid string) *ImageInsertedDTO {
	return &ImageInsertedDTO{
		ImageUid: imageUid,
	}
}
