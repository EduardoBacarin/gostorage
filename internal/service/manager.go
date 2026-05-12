package service

import (
	"github.com/EduardoBacarin/gostorage/internal/security"
	"github.com/EduardoBacarin/gostorage/internal/storage"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Services struct {
	Object  *ObjectService
	User    *UserService
	Session *security.SessionManager
}

func NewServices(db *mongo.Database) *Services {
	sessionManager := security.NewSessionManager()
	store := storage.NewLocalStorage("./store")
	return &Services{
		Session: sessionManager,
		User:    NewUserService(db, sessionManager),
		Object:  NewObjectService(store, db),
	}
}
