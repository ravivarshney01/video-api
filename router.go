package main

import (
	"github.com/gorilla/mux"
	"net/http"
	handler "video-api/handler"
	"video-api/response"
)

func addRoutes(r *mux.Router) {
	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)

	h := handler.Handler{}
	r.HandleFunc("/upload", h.UploadVideo).Methods(http.MethodPost)
	//r.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
	//	handleVideoUpload(w, r, db)
	//}).Methods("POST")

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response.WithJSON(w, http.StatusOK, "ok", nil)
	return
}
