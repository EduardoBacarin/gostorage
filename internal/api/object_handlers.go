package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (h *Handler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")
	object := r.PathValue("object")
	if bucket == "" || object == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Invalid file or bucket")
		return
	}

	buffer := make([]byte, 512)
	n, err := r.Body.Read(buffer)
	if err != nil && err != io.EOF {
		SendJSON(w, http.StatusInternalServerError, false, nil, "Error reading upload stream")
		return
	}

	contentType := "application/octet-stream"
	if n > 0 {
		contentType = http.DetectContentType(buffer[:n])
	}
	fullStream := io.MultiReader(bytes.NewReader(buffer[:n]), r.Body)

	metadata, err := h.srv.Object.Upload(r.Context(), fullStream, bucket, object, contentType)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	data := map[string]string{"id": metadata.ID}
	SendJSON(w, http.StatusCreated, true, data, "")
}

func (h *Handler) GetObjectHandler(w http.ResponseWriter, r *http.Request) {
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
