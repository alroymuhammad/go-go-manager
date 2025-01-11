package route

import (
	"database/sql"

	"github.com/alroymuhammad/go-go-manager/internal/handlers"
	"github.com/alroymuhammad/go-go-manager/internal/middleware"
	"github.com/gorilla/mux"
)

// SetupRoutes mengatur routing untuk aplikasi.
func SetupRoutes(db *sql.DB) *mux.Router {
	// Inisialisasi handler untuk departments
	departmentHandler := handlers.NewDepartmentsHandler(db)
	employeeHandler := handlers.NewEmployeesHandler(db)

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

	// Kamu bisa menambahkan middleware jika perlu di sini
	router.Use(middleware.LoggingMiddleware)

	return router
}
