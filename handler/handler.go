package handler

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"video-api/constants"
	"video-api/models"
	"video-api/request"
	"video-api/response"
	"video-api/uploadSvc"
)

type Handler struct{}

func (h *Handler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(constants.MaximumFileToUpload)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	file, header, err := r.FormFile("video")
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer file.Close()

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".mp4") {
		response.WithError(w, http.StatusUnsupportedMediaType, "only mp4 is supported")
		return
	}

	if header.Size > constants.MaximumFileToUpload {
		http.Error(w, "File size exceeds maximum limit", http.StatusBadRequest)
		return
	}

	uploadService := uploadSvc.NewFileSystemUploadService()
	path, err := uploadService.Upload(header.Filename, file)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := models.AddVideo(r.Context(), header.Filename, path)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := map[string]interface{}{
		"id": id,
	}

	response.WithJSON(w, http.StatusOK, "uploaded", data)
}

func (h *Handler) TrimVideo(w http.ResponseWriter, r *http.Request) {

	start := request.ParseIntQueryParam(r, "start", 0)
	end := request.ParseIntQueryParam(r, "end", 0)

	if start >= end {
		response.WithError(w, http.StatusInternalServerError, "start must be less than end")
		return
	}

	videoIdStr := r.FormValue("id")
	videoId, err := strconv.Atoi(videoIdStr)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	video, err := models.GetVideo(r.Context(), videoId)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	outputPath := fmt.Sprintf("uploads/trimmed_%s", video.Filename)
	cmd := exec.Command("ffmpeg", "-ss", fmt.Sprintf("%d", start), "-to", fmt.Sprintf("%d", end), "-i", video.Url, "-c", "copy", outputPath)
	if err := cmd.Run(); err != nil {
		response.WithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	video.Url = outputPath
	video.Filename = "trimmed_" + video.Filename
	err = video.UpdateVideo(r.Context())
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, "unable to update video path "+err.Error())
		return
	}
	response.WithJSON(w, http.StatusOK, "trimmed successfully", nil)
	return
}

func (h *Handler) MergeVideos(w http.ResponseWriter, r *http.Request) {

}
