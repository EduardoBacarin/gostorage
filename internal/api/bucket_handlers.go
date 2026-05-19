package api

import (
	"encoding/json"
	"net/http"

	"github.com/EduardoBacarin/gostorage/internal/security"
)

type CreateBucketRequest struct {
	Name   string `json:"name"`
	Public bool   `json:"public"`
}

func (h *Handler) CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(SessionKey).(security.SessionData)
	if !ok {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	var req CreateBucketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSON(w, http.StatusBadRequest, false, nil, "Invalid JSON payload")
		return
	}

	bucket, err := h.srv.Bucket.CreateBucket(r.Context(), req.Name, req.Public, session.UserID)
	if err != nil {
		if err.Error() == "Bucket already exists" {
			SendJSON(w, http.StatusConflict, false, nil, err.Error())
			return
		}

		SendJSON(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	SendJSON(w, http.StatusCreated, true, bucket, "")
}

type UpdateBucketRequest struct {
	Name   *string `json:"name,omitempty"`
	Public *bool   `json:"public,omitempty"`
}

func (h *Handler) UpdateBucketHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(SessionKey).(security.SessionData)
	if !ok {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	bucketID := r.PathValue("id")
	if bucketID == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Invalid bucket ID")
		return
	}

	var req UpdateBucketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSON(w, http.StatusBadRequest, false, nil, "Invalid JSON payload")
		return
	}

	updatedBucket, err := h.srv.Bucket.UpdateBucket(r.Context(), bucketID, req.Name, req.Public, session.UserID)
	if err != nil {
		switch err.Error() {
		case "Not Found":
			SendJSON(w, http.StatusNotFound, false, nil, err.Error())
		case "Permission denied":
			SendJSON(w, http.StatusForbidden, false, nil, "You do not have permission to modify this bucket")
		case "Bucket already exists":
			SendJSON(w, http.StatusConflict, false, nil, err.Error())
		case "Bucket name cannot be empty":
			SendJSON(w, http.StatusBadRequest, false, nil, err.Error())
		default:
			SendJSON(w, http.StatusInternalServerError, false, nil, err.Error())
		}
		return
	}

	SendJSON(w, http.StatusOK, true, updatedBucket, "")
}

func (h *Handler) GetBucketHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(SessionKey).(security.SessionData)
	if !ok {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	bucketID := r.PathValue("id")
	if bucketID == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Invalid bucket ID")
		return
	}

	bucket, err := h.srv.Bucket.GetBucket(r.Context(), bucketID, &session.UserID)
	if err != nil {
		switch err.Error() {
		case "Not Found":
			SendJSON(w, http.StatusNotFound, false, nil, err.Error())
		default:
			SendJSON(w, http.StatusInternalServerError, false, nil, err.Error())
		}
		return
	}

	SendJSON(w, http.StatusOK, true, bucket, "")
}

func (h *Handler) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(SessionKey).(security.SessionData)
	if !ok {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	bucketID := r.PathValue("id")
	if bucketID == "" {
		SendJSON(w, http.StatusBadRequest, false, nil, "Invalid bucket ID")
		return
	}

	err := h.srv.Bucket.DeleteBucket(r.Context(), bucketID, &session.UserID)
	if err != nil {
		switch err.Error() {
		case "Not Found":
			SendJSON(w, http.StatusNotFound, false, nil, err.Error())
		default:
			SendJSON(w, http.StatusInternalServerError, false, nil, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
