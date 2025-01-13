package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"log"
)

type EmployeesHandler struct {
	db *sql.DB
}

type EmployeesRequest struct {
	IdentityNumber   string `json:"identityNumber"`
	Name             string `json:"name"`
	EmployeeImageUri string `json:"employeeImageUri"`
	Gender           string `json:"gender"`
	DepartmentID     int    `json:"department_id"`
	ManagerID        int    `json:"manager_id"`
}

type EmployeesResponse struct {
	IdentityNumber   string    `json:"identityNumber"`
	Name             string    `json:"name"`
	EmployeeImageUri string    `json:"employeeImageUri"`
	Gender           string    `json:"gender"`
	DepartmentID     int       `json:"departmentId"`
}

func NewEmployeesHandler(db *sql.DB) *EmployeesHandler {
	return &EmployeesHandler{db: db}
}

func (h *EmployeesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/employees":
		h.CreateEmployees(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/employees":
		h.ListEmployees(w, r)
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
		"INSERT INTO employees (identityNumber, name, employeeImageUri, gender, department_id, manager_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING identityNumber",
		req.IdentityNumber, req.Name, req.EmployeeImageUri, req.Gender, req.DepartmentID, req.ManagerID,
	).Scan(&identity_number)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"identity_number": identity_number})
}


// Get Employee
func (h *EmployeesHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	// Set default values if limit and offset are not provided
	if limit == "" {
		limit = "10" // Default limit
	}
	if offset == "" {
		offset = "0" // Default offset
	}

	// Query to get employees
	rows, err := h.db.Query(
		"SELECT identityNumber, name, employeeImageUri, gender, department_id FROM employees LIMIT $1 OFFSET $2",
		limit, offset,
	)

	log.Printf("Request Data: %+v\n", rows)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()



	var employees []EmployeesResponse
	for rows.Next() {
		var emp EmployeesResponse
		err := rows.Scan(
			&emp.IdentityNumber,
			&emp.Name,
			&emp.EmployeeImageUri,
			&emp.Gender,
			&emp.DepartmentID,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		employees = append(employees, emp)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get total count of employees for pagination info
	var totalCount int
	err = h.db.QueryRow("SELECT COUNT(*) FROM employees").Scan(&totalCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create response object with pagination data
	response := struct {
		TotalCount int                `json:"totalCount"`
		Limit      string             `json:"limit"`
		Offset     string             `json:"offset"`
		Employees  []EmployeesResponse `json:"employees"`
	}{
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
		Employees:  employees,
	}

	// Send the response as JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

