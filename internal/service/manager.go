package service

import (
	"github.com/EduardoBacarin/gostorage/internal/storage"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Services struct {
	Object *ObjectService
}

func NewServices(db *mongo.Database, store storage.Engine) *Services {
	return &Services{
		Object: NewObjectService(store, db),
	}
}
