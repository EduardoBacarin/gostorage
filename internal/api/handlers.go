package api

import (
	"github.com/EduardoBacarin/gostorage/internal/service"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Handler struct {
	db  *mongo.Database
	srv *service.Services
}

func NewHandler(db *mongo.Database, services *service.Services) *Handler {
	return &Handler{
		db:  db,
		srv: services,
	}
}
