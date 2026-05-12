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
	client, _ := database.ConnectMongo("mongodb://localhost:27017")
	db := client.Database("gostorage")

	store := storage.NewLocalStorage("./store")
	services := service.NewServices(db, store)

	h := api.NewHandler(db, services)

	router := api.SetupRoutes(h)

	log.Println("Server started at port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
