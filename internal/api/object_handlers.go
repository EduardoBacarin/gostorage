package api

import (
	"bytes"
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

func (h *Handler) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.PathValue("bucket")
	id := r.PathValue("id")

	if bucket == "" || id == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Invalid bucket or Id")
		return
	}

	metadata, file, err := h.srv.Object.GetObject(r.Context(), bucket, id)
	if err != nil {
		if err.Error() == "Not found" {
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(500)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", metadata.ContentType)

	w.Header().Set("Content-Length", strconv.FormatInt(metadata.Size, 10))
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("Erro ao transmitir arquivo para o cliente: %v", err)
	}
}
