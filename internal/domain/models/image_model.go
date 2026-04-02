package models

type ImageModel struct {
	UID         string
	FileName    string
	Extension   string
	FullPath    string
	Object      []byte
	Status      Status
	ProcessType ProcessType
}

func NewImageModel(uid, fileName, extension, fullPath string, object []byte, processType ProcessType, status Status) *ImageModel {
	return &ImageModel{
		UID:         uid,
		FileName:    fileName,
		Extension:   extension,
		FullPath:    fullPath,
		Object:      object,
		Status:      status,
		ProcessType: processType,
	}
}
