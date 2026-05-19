package security

import (
	"sync"
	"time"

	"github.com/EduardoBacarin/gostorage/internal/helpers"
)

type SessionData struct {
	UserID      string    `json:"id"`
	Email       string    `json:"email"`
	Permissions []string  `json:"permissions"`
	Buckets     []string  `json:"buckets"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]SessionData
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]SessionData),
	}
}

func (m *SessionManager) CreateSession(userID string, email string, permissions []string, buckets []string, duration time.Duration) string {
	token := helpers.GenerateSHA256(email, time.Now().String(), "session-secret")

	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions[token] = SessionData{
		UserID:      userID,
		Email:       email,
		Permissions: permissions,
		Buckets:     buckets,
		ExpiresAt:   time.Now().Add(duration),
	}

	return token
}

func (m *SessionManager) GetSession(token string) (SessionData, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, exists := m.sessions[token]
	if !exists {
		return SessionData{}, false
	}

	if time.Now().After(data.ExpiresAt) {
		return SessionData{}, false
	}

	return data, true
}

func (s *SessionData) HasPermission(required string) bool {
	for _, p := range s.Permissions {
		if p == required {
			return true
		}
	}
	return false
}

func (s *SessionData) CanAccessBucket(bucketID string) bool {
	for _, b := range s.Buckets {
		if b == "*" || b == bucketID {
			return true
		}
	}
	return false
}

func (s *SessionManager) DeleteSession(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
}
