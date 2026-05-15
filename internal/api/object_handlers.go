package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/EduardoBacarin/gostorage/internal/helpers"
	"github.com/EduardoBacarin/gostorage/internal/security"
)

func (h *Handler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value(SessionKey).(security.SessionData)
	bucket := r.PathValue("bucket")
	object := r.PathValue("object")
	log.Println(object, bucket)
	if bucket == "" || object == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Invalid file or bucket")
		return
	}

	maxFileSize, _ := strconv.ParseInt(os.Getenv("MAX_FILE_SIZE"), 10, 64)
	if r.ContentLength > maxFileSize {
		SendJSON(w, http.StatusRequestEntityTooLarge, false, nil, "Content exceeds the maximum allowed file size")
		return
	}

	findBucket, err := h.srv.Bucket.GetBucket(r.Context(), helpers.GenerateSHA256("bucket", bucket), nil)
	if err != nil {
		SendJSON(w, http.StatusNotFound, false, nil, "Bucket not found")
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

	metadata, err := h.srv.Object.Upload(r.Context(), fullStream, findBucket.Name, object, session.UserID, contentType, r.ContentLength)
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
