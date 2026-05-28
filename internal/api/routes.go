package api

import "net/http"

func SetupRoutes(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	/* AUTH */
	mux.HandleFunc("POST /v1/auth", h.LoginHandler)
	mux.Handle("GET /v1/auth", h.AuthMiddleware(http.HandlerFunc(h.MeHandler)))
	mux.Handle("DELETE /v1/auth", h.AuthMiddleware(http.HandlerFunc(h.LogoutHandler)))

	/* BUCKET */
	mux.Handle("POST /v1/bucket", h.AuthMiddleware(http.HandlerFunc(h.CreateBucketHandler)))
	mux.Handle("PATCH /v1/bucket/{id}", h.AuthMiddleware(http.HandlerFunc(h.UpdateBucketHandler)))
	mux.Handle("GET /v1/bucket/{id}", h.AuthMiddleware(http.HandlerFunc(h.GetBucketHandler)))
	mux.Handle("GET /v1/bucket", h.AuthMiddleware(http.HandlerFunc(h.ListBucketHandler)))
	mux.Handle("DELETE /v1/bucket/{id}", h.AuthMiddleware(http.HandlerFunc(h.DeleteBucketHandler)))

	/* OBJECT */
	mux.Handle("PUT /{bucket}/{object}", h.AuthMiddleware(http.HandlerFunc(h.UploadObjectHandler)))
	mux.Handle("DELETE /{bucket}/{id}", h.AuthMiddleware(http.HandlerFunc(h.DeleteObjectHandler)))
	mux.HandleFunc("GET /{bucket}/{id}", h.DownloadObjectHandler)
	/* mux.HandleFunc("DELETE /v1/object/{id}", h.DeleteHandler) */

	return mux
}
