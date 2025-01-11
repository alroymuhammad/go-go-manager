package main

import (
	"fmt"
	"log"
	"net/http"

	route "github.com/alroymuhammad/go-go-manager/internal/routes"
	config "github.com/alroymuhammad/go-go-manager/pkg/database"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func main() {
	// Connect to the database
	db := config.ConnectDB()
	defer db.Close()

	// Set up routes
	router := route.SetupRoutes(db)
	router.HandleFunc("/", home)

	// Create a new HTTP server and pass the router
	port := 8080
	fmt.Printf("Starting server on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
