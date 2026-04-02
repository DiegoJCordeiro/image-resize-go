package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"

	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/infra/http/handlers/handlersdto"
	"github.com/DiegoJCordeiro/image-resizing-go-api/internal/usecases"
)

type ImageHandler struct {
	getImageUseCase    usecases.GetImageUseCase
	insertImageUseCase usecases.InsertImageUseCase
	deleteImageUseCase usecases.DeleteImageUseCase
}

func NewImageHandler(
	getImageUseCase usecases.GetImageUseCase,
	insertImageUseCase usecases.InsertImageUseCase,
	deleteImageUseCase usecases.DeleteImageUseCase,
) *ImageHandler {
	return &ImageHandler{
		getImageUseCase:    getImageUseCase,
		insertImageUseCase: insertImageUseCase,
		deleteImageUseCase: deleteImageUseCase,
	}
}

func (iih *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {

	imageUid := r.PathValue("image-uid")

	imageModel, err := iih.getImageUseCase.Execute(imageUid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	imageDto := handlersdto.NewImageDTO(
		imageModel.UID,
		imageModel.Object,
	)

	imageMarshal, err := json.Marshal(imageDto)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(imageMarshal)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (iih *ImageHandler) InsertImage(w http.ResponseWriter, r *http.Request) {

	var fileInBytes []byte
	formFile, headerFile, err := r.FormFile("imageFile")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer formFile.Close()

	fileInBytes, err = io.ReadAll(formFile)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = iih.insertImageUseCase.Execute(
		headerFile.Filename,
		filepath.Ext(headerFile.Filename),
		r.Header.Get("x-process-type"),
		fileInBytes,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}

func (iih *ImageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {

	imageUid := r.PathValue("image-uid")

	err := iih.deleteImageUseCase.Execute(imageUid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (iih *ImageHandler) GetAllHandlers() map[string]func(http.ResponseWriter, *http.Request) {

	handlersMap := map[string]func(http.ResponseWriter, *http.Request){
		"POST /image":               iih.InsertImage,
		"GET /image/{image-uid}":    iih.GetImage,
		"DELETE /image/{image-uid}": iih.DeleteImage,
	}

	return handlersMap
}
