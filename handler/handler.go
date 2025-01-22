package handler

import (
	"net/http"
	"strconv"
	"strings"
	"video-api/constants"
	"video-api/core"
	"video-api/request"
	"video-api/response"
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

	id, err := core.UploadVideo(r.Context(), file, header.Filename)
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

	err = core.TrimVideo(r.Context(), videoId, start, end)
	response.WithJSON(w, http.StatusOK, "trimmed successfully", nil)
	return
}

func (h *Handler) MergeVideos(w http.ResponseWriter, r *http.Request) {
	videoIds := request.ParseCommaSeparatedQueryParamIds(r, "ids")
	if len(videoIds) == 0 {
		response.WithError(w, http.StatusInternalServerError, "ids is required")
		return
	}
	id, err := core.MergeVideos(r.Context(), videoIds)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WithJSON(w, http.StatusOK, "merged successfully", map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) ShareVideo(w http.ResponseWriter, r *http.Request) {
	videoIdStr := r.FormValue("id")
	videoId, err := strconv.Atoi(videoIdStr)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	url, err := core.ShareVideoUrl(r.Context(), videoId)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WithJSON(w, http.StatusOK, "shared successfully", map[string]interface{}{
		"url": url,
	})
}

func (h *Handler) GetVideo(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	videoLink, err := core.ValidateJWT(token)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Content-Disposition", "inline")
	http.ServeFile(w, r, videoLink)
	return
}
