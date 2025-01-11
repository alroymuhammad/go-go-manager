package models

import "time"

type Employees struct {
	IdentityNumber   string    `json:"identity_number"`
	Name             string    `json:"name"`
	EmployeeImageUri string    `json:"employee_image_uri"`
	Gender           string    `json:"gender"`
	DepartmentID     int       `json:"department_id"`
	ManagerID        int       `json:"manager_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
