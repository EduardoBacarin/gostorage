package api

import "net/http"

func SetupRoutes(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", h.UploadHandler)
	return mux
}
