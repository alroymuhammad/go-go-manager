package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alroymuhammad/go-go-manager/internal/routes"
	config "github.com/alroymuhammad/go-go-manager/pkg/database"
)

func main() {
	db := config.ConnectDB()
	defer db.Close()

	router := routes.NewRouter(db)

	port := 8080
	fmt.Printf("Starting server on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	log.Fatal(err)
}
