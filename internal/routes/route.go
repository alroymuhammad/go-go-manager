package routes

import (
	"database/sql"

	"github.com/alroymuhammad/go-go-manager/internal/handlers"
	"github.com/alroymuhammad/go-go-manager/internal/handlers/auth_handler"
	"github.com/alroymuhammad/go-go-manager/internal/middleware"
	"github.com/alroymuhammad/go-go-manager/internal/usecase/auth_usecase"
	"github.com/gorilla/mux"
)

// SetupRoutes mengatur routing untuk aplikasi.
func SetupRoutes(db *sql.DB) *mux.Router {
	// Inisialisasi handler untuk departments
	departmentHandler := handlers.NewDepartmentsHandler(db)
	employeeHandler := handlers.NewEmployeesHandler(db)
	authService := auth_usecase.NewAuthService(db)
	authHandler := auth_handler.NewAuthHandler(authService)

	// Inisialisasi router mux
	router := mux.NewRouter()

	// Routing untuk Departments
	router.Handle("/departments", departmentHandler).Methods("POST")        // Create Department
	router.Handle("/departments", departmentHandler).Methods("GET")         // List Departments
	router.Handle("/departments/{id}", departmentHandler).Methods("GET")    // Get Department by ID
	router.Handle("/departments/{id}", departmentHandler).Methods("PUT")    // Update Department by ID
	router.Handle("/departments/{id}", departmentHandler).Methods("DELETE") // Delete Department by ID

	// Routing untuk Employees
	router.Handle("/employees", employeeHandler).Methods("POST") // Create Employee

	// Routing untuk Auth
	router.HandleFunc("/v1/auth", authHandler.AuthHandler).Methods("POST")

	// Kamu bisa menambahkan middleware jika perlu di sini
	router.Use(middleware.LoggingMiddleware)

	return router
}
