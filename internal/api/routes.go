package api

import "net/http"

func SetupRoutes(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", h.UploadHandler)
	mux.HandleFunc("GET /{id}", h.RetrieveHandler)
	mux.HandleFunc("DELETE /{id}", h.DeleteHandler)
	return mux
}
