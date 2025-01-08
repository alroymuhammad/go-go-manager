package handlers/auth
import (
	"github.com/alroymuhammad/go-go-manager/auth/jwt"
	"github.com/alroymuhammad/go-go-manager/internal/models/users"
)

func Login(email, password string) (string, error) {
	user, err := users.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	return jwt.GenerateJWT(user.ID), nil
}