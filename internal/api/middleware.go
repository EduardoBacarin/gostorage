package api

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const SessionKey ctxKey = "session_data"

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
			return
		}

		session, ok := h.srv.Session.GetSession(token)
		if !ok {
			SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), SessionKey, session)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
