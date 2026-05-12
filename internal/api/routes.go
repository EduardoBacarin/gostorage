package api

import "net/http"

func SetupRoutes(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	/* AUTH */
	mux.HandleFunc("POST /v1/auth", h.LoginHandler)
	mux.Handle("GET /v1/auth", h.AuthMiddleware(http.HandlerFunc(h.MeHandler)))
	mux.Handle("DELETE /v1/auth", h.AuthMiddleware(http.HandlerFunc(h.LogoutHandler)))

	/* OBJECT */
	mux.HandleFunc("POST /v1/object", h.UploadHandler)
	mux.HandleFunc("GET /v1/object/{id}", h.RetrieveHandler)
	mux.HandleFunc("DELETE /v1/object/{id}", h.DeleteHandler)
	return mux
}
