package handlersdto

type ImageDTO struct {
	ImageUid    string `json:"image_uid"`
	ImageObject []byte `json:"image_object"`
}

func NewImageDTO(imageUid string, imageObject []byte) *ImageDTO {
	return &ImageDTO{
		ImageUid:    imageUid,
		ImageObject: imageObject,
	}
}
