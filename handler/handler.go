package handler

import (
	"fmt"
	"net/http"
	"video-api/response"
	"video-api/uploadSvc"
)

type Handler struct{}

func (h *Handler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Failed to read uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadService := uploadSvc.NewFileSystemUploadService()
	path, err := uploadService.Upload(header.Filename, file)
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusBadRequest)
		return
	}

	fmt.Println(path)

	response.WithJSON(w, http.StatusOK, "uploaded", nil)
}
