package routes

import (
	"database/sql"

	"github.com/alroymuhammad/go-go-manager/internal/handlers/auth_handler"
	"github.com/alroymuhammad/go-go-manager/internal/usecase/auth_usecase"
	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	authService := auth_usecase.NewAuthService(db)
	authHandler := auth_handler.NewAuthHandler(authService)

	router.HandleFunc("/v1/auth", authHandler.AuthHandler).Methods("POST")

	return router
}
