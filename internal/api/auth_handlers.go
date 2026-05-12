package api

import (
	"encoding/json"
	"net/http"

	"github.com/EduardoBacarin/gostorage/internal/security"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		SendJSON(w, http.StatusBadRequest, false, nil, "Bad Request")
		return
	}

	token, err := h.srv.User.Authenticate(r.Context(), input.Email, input.Password)
	if err != nil {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Invalid Credentials")
		return
	}
	var data = map[string]string{
		"token": token,
	}
	SendJSON(w, http.StatusOK, true, data, "")
}

func (h *Handler) MeHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(SessionKey).(security.SessionData)
	if !ok {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}
	SendJSON(w, http.StatusOK, true, session, "")
}
