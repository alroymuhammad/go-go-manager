package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type DepartmentsHandler struct {
	db *sql.DB
}

type DepartmentsRequest struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type DepartmentsResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewDepartmentsHandler(db *sql.DB) *DepartmentsHandler {
	return &DepartmentsHandler{db: db}
}

func getDepartmentsIDFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		return parts[2]
	}
	return ""
}

func (h *DepartmentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/departments":
		h.CreateDepartments(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/departments":
		h.ListDepartments(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/departments/"):
		h.GetDepartments(w, r)
	case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/departments/"):
		h.UpdateDepartments(w, r)
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/departments/"):
		h.DeleteDepartments(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *DepartmentsHandler) CreateDepartments(w http.ResponseWriter, r *http.Request) {
	var req DepartmentsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var id int64
	err := h.db.QueryRow(
		"INSERT INTO departments (id, name, created_at, updated_at) VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id",
		req.ID, req.Name,
	).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func (h *DepartmentsHandler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, name, created_at, updated_at FROM departments")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var departments []DepartmentsResponse
	for rows.Next() {
		var department DepartmentsResponse
		if err := rows.Scan(&department.ID, &department.Name, &department.CreatedAt, &department.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		departments = append(departments, department)
	}

	json.NewEncoder(w).Encode(departments)
}

func (h *DepartmentsHandler) GetDepartments(w http.ResponseWriter, r *http.Request) {
	id := getDepartmentsIDFromPath(r.URL.Path)

	var department DepartmentsResponse
	err := h.db.QueryRow("SELECT id, name, created_at, updated_at FROM departments WHERE id = $1", id).Scan(
		&department.ID, &department.Name, &department.CreatedAt, &department.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(department)
}

func (h *DepartmentsHandler) UpdateDepartments(w http.ResponseWriter, r *http.Request) {
	id := getDepartmentsIDFromPath(r.URL.Path)

	var req DepartmentsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := h.db.Exec("UPDATE departments SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		req.Name, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DepartmentsHandler) DeleteDepartments(w http.ResponseWriter, r *http.Request) {
	id := getDepartmentsIDFromPath(r.URL.Path)

	_, err := h.db.Exec("DELETE FROM departments WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
