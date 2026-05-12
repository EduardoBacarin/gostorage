package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/EduardoBacarin/gostorage/internal/api"
	"github.com/EduardoBacarin/gostorage/internal/database"
	"github.com/EduardoBacarin/gostorage/internal/service"
)

func main() {
	setupMode := flag.Bool("setup", false, "Execute the first start setup")
	flag.Parse()
	client, _ := database.ConnectMongo("mongodb://localhost:27017")
	db := client.Database("gostorage")

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

	log.Println("Server started at port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
