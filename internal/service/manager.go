package service

import (
	"os"

	"github.com/EduardoBacarin/gostorage/internal/security"
	"github.com/EduardoBacarin/gostorage/internal/storage"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Services struct {
	Object  *ObjectService
	User    *UserService
	Session *security.SessionManager
	Bucket  *BucketService
}

func NewServices(db *mongo.Database) *Services {
	sessionManager := security.NewSessionManager()
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage" // Fallback de segurança caso a variável falhe
	}
	store := storage.NewLocalStorage(storagePath)
	return &Services{
		Session: sessionManager,
		User:    NewUserService(db, sessionManager),
		Object:  NewObjectService(store, db),
		Bucket:  NewBucketService(db),
	}
}
