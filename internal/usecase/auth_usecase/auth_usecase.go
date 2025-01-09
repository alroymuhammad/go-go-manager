package auth_usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	jwt "github.com/alroymuhammad/go-go-manager/internal/appjwt"
	models "github.com/alroymuhammad/go-go-manager/internal/models/users"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidInput       = errors.New("invalid input parameters")
	ErrUserExists         = errors.New("user already exists")
)

type AuthService struct {
	DB *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) ValidateCredentials(email, password string) error {
	if !isValidEmail(email) {
		return fmt.Errorf("invalid email format")
	}
	if !isValidPassword(password) {
		return fmt.Errorf("password must be between 8 and 32 characters")
	}
	return nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	if err := s.ValidateCredentials(email, password); err != nil {
		return "", ErrInvalidInput
	}

	var user models.User
	err := s.DB.QueryRow("SELECT id, email, password FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return jwt.GenerateJWT(fmt.Sprintf("%d", user.ID))
}

func (s *AuthService) Register(email, password string) (*models.User, error) {
	if err := s.ValidateCredentials(email, password); err != nil {
		return nil, ErrInvalidInput
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.DB.Exec("INSERT INTO users (email, password, created_at, updated_at) VALUES ($1, $2, $3, $4)",
		user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, ErrUserExists
		}
		return nil, err
	}

	return user, nil
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func isValidPassword(password string) bool {
	return len(password) >= 8 && len(password) <= 32
}
