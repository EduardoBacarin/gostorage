package api

import (
	"encoding/json"
	"net/http"

	"github.com/EduardoBacarin/gostorage/internal/service"
)

type Handler struct {
	objService *service.ObjectService
}

func NewHandler(s *service.ObjectService) *Handler {
	return &Handler{objService: s}
}

func (h *Handler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File not found", http.StatusBadRequest)
		return
	}
	defer file.Close()

	bucket := r.FormValue("bucket")
	if bucket == "" {
		bucket = "default"
	}
	key := header.Filename

	metadata, err := h.objService.Upload(r.Context(), file, bucket, key)
	if err != nil {
		http.Error(w, "Erro ao processar upload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(metadata)
}
