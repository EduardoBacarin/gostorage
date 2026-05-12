package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func (h *Handler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		SendJSON(w, http.StatusBadRequest, false, nil, "Invalid file")
		return
	}
	defer file.Close()

	bucket := r.FormValue("bucket")
	if bucket == "" {
		bucket = "default"
	}
	key := header.Filename

	metadata, err := h.srv.Object.Upload(r.Context(), file, bucket, key)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	data := map[string]string{"id": metadata.ID}
	SendJSON(w, http.StatusCreated, true, data, "")
}

func (h *Handler) RetrieveHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	stream, meta, err := h.srv.Object.GetObject(r.Context(), id)
	if err != nil {
		SendJSON(w, http.StatusNotFound, false, nil, "")
		return
	}
	defer stream.Close()

	w.Header().Set("Content-Type", meta.ContentType)

	if meta.Size > 0 {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", meta.Size))
	}

	n, err := io.Copy(w, stream)
	if err != nil {
		log.Printf("Streaming error: %v", err)
	}
	log.Printf("Sent %d bytes", n)
}

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "")
		return
	}
	err := h.srv.Object.DeleteObject(r.Context(), id)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	SendJSON(w, http.StatusOK, true, nil, "")
}
