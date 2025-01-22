package main

import (
	"github.com/gorilla/mux"
	"net/http"
	handler "video-api/handler"
	"video-api/response"
)

func addRoutes(publicRouter, authenticatedRouter *mux.Router) {
	publicRouter.HandleFunc("/health", healthHandler).Methods(http.MethodGet)

	h := handler.Handler{}
	authenticatedRouter.HandleFunc("/upload", h.UploadVideo).Methods(http.MethodPost)
	authenticatedRouter.HandleFunc("/trim", h.TrimVideo).Methods(http.MethodPost)
	authenticatedRouter.HandleFunc("/merge", h.MergeVideos).Methods(http.MethodPost)
	authenticatedRouter.HandleFunc("/share-video", h.ShareVideo).Methods(http.MethodGet)
	publicRouter.HandleFunc("/video", h.GetVideo).Methods(http.MethodGet)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response.WithJSON(w, http.StatusOK, "ok", nil)
	return
}
