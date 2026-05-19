package service

import (
	"context"
	"errors"
	"time"

	"github.com/EduardoBacarin/gostorage/internal/helpers"
	"github.com/EduardoBacarin/gostorage/internal/models"
	"github.com/EduardoBacarin/gostorage/internal/security"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserService struct {
	userCol  *mongo.Collection
	sessions *security.SessionManager
}

func NewUserService(db *mongo.Database, sm *security.SessionManager) *UserService {
	return &UserService{
		userCol:  db.Collection("users"),
		sessions: sm,
	}
}

func (s *UserService) CreateUser(ctx context.Context, email string, password string, permissions []string, buckets []string) error {
	var existingUser models.User
	err := s.userCol.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		return errors.New("User already registered")
	}

	userID := helpers.GenerateSHA256("user", email, helpers.NowStrNano())
	passwordHash := helpers.GenerateSHA256("gostorage", "pass", password)

	user := models.User{
		ID:          userID,
		Email:       email,
		Password:    passwordHash,
		Permissions: permissions,
		Buckets:     buckets,
		CreatedAt:   time.Now(),
	}

	_, err = s.userCol.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Authenticate(ctx context.Context, email, password string) (string, error) {
	sentPasswordHash := helpers.GenerateSHA256("gostorage", "pass", password)

	var user models.User
	err := s.userCol.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return "", errors.New("Invalid Credentials")
	}

	if user.Password != sentPasswordHash {
		return "", errors.New("Invalid Credentials")
	}

	token := s.sessions.CreateSession(
		user.ID,
		user.Email,
		user.Permissions,
		user.Buckets,
		2*time.Hour,
	)

	return token, nil
}
