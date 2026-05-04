package main

import (
	"log"
	"net/http"

	"github.com/EduardoBacarin/gostorage/internal/api"
	"github.com/EduardoBacarin/gostorage/internal/database"
	"github.com/EduardoBacarin/gostorage/internal/service"
	"github.com/EduardoBacarin/gostorage/internal/storage"
)

func main() {
	client, err := database.ConnectMongo("mongodb://localhost:27017")
	if err != nil {
		log.Fatal("Falha ao conectar no Mongo:", err)
	}
	collection := client.Database("gostorage").Collection("objects")

	store := storage.NewLocalStorage("./store")

	objService := service.NewObjectService(store, collection)

	h := api.NewHandler(objService)

	router := api.SetupRoutes(h)

	log.Println("Server started at port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
