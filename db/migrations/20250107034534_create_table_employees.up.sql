CREATE TABLE employees (
    identityNumber VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    employeeImageUri VARCHAR(255),
    gender VARCHAR(255),
    department_id INTEGER NOT NULL,
    manager_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (department_id) REFERENCES departments(id),
    FOREIGN KEY (manager_id) REFERENCES users(id)
);