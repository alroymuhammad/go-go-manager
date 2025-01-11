package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type EmployeesHandler struct {
	db *sql.DB
}

type EmployeesRequest struct {
	IdentityNumber   string `json:"identity_number"`
	Name             string `json:"name"`
	EmployeeImageUri string `json:"employee_image_uri"`
	Gender           string `json:"gender"`
	DepartmentID     int    `json:"department_id"`
	ManagerID        int    `json:"manager_id"`
}

type EmployeesResponse struct {
	IdentityNumber   string    `json:"identity_number"`
	Name             string    `json:"name"`
	EmployeeImageUri string    `json:"employee_image_uri"`
	Gender           string    `json:"gender"`
	DepartmentID     int       `json:"department_id"`
	ManagerID        int       `json:"manager_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func NewEmployeesHandler(db *sql.DB) *EmployeesHandler {
	return &EmployeesHandler{db: db}
}

func (h *EmployeesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/employees":
		h.CreateEmployees(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *EmployeesHandler) CreateEmployees(w http.ResponseWriter, r *http.Request) {
	var req EmployeesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var identity_number int64
	err := h.db.QueryRow(
		"INSERT INTO employees (identity_number, name, employee_image_uri, gender, department_id, manager_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING identity_number",
		req.IdentityNumber, req.Name, req.EmployeeImageUri, req.Gender, req.DepartmentID, req.ManagerID,
	).Scan(&identity_number)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"identity_number": identity_number})
}
