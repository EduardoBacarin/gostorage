package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/EduardoBacarin/gostorage/internal/api"
	"github.com/EduardoBacarin/gostorage/internal/database"
	"github.com/EduardoBacarin/gostorage/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	setupMode := flag.Bool("setup", false, "Execute the first start setup")
	flag.Parse()
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env File not Found")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Valor padrão caso não esteja no .env
	}
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")
	if mongoURI == "" || dbName == "" {
		log.Fatal("MONGO_URI and MONGO_DB_NAME are required")
	}

	client, err := database.ConnectMongo(mongoURI)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}
	db := client.Database(dbName)

	services := service.NewServices(db)

	if *setupMode {
		fmt.Println("🛠️  Welcome to GoStorage Setup")
		fmt.Println("-----------------------------------")

		var email, password string

		fmt.Print("📧 Type admin e-mail: ")
		fmt.Scanln(&email)

		fmt.Print("🔑 Type admin password: ")
		fmt.Scanln(&password)

		fmt.Printf("\nCreating admin: %s...\n", email)

		err := services.User.CreateUser(
			context.Background(),
			email,
			password,
			[]string{"*"},
			[]string{"*"},
		)

		if err != nil {
			log.Fatalf("ERROR: %v", err)
		}

		fmt.Println("\n✅ Setup finished succesfully, welcome aboard!")
		return
	}

	h := api.NewHandler(db, services)
	router := api.SetupRoutes(h)

	log.Printf("Server started at port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
