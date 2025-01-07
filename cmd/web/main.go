package main

import (
	"fmt"
	"log"
	"net/http"
	
	"github.com/alroymuhammad/go-go-manager/pkg/database"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func main() {
	db := config.ConnectDB()
	defer db.Close()
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	port := 8080
	fmt.Printf("Starting server on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	log.Fatal(err)
}
