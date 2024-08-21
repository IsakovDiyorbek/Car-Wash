package token

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Car-Wash/Car-Wash-Api-Gateway/config"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
)

type JWTHandler struct {
	Sub        string
	Exp        string
	Iat        string
	Role       string
	SigningKey string
	Token      string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

var tokenKey = config.Load().TokenKey

// CreateToken creates a new token

func ExtractClaim(tokenStr string) (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey), nil
	}
	token, err = jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ExtractClaims ...
func ExtractClaims(jwttoken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(jwttoken, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Load().TokenKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		slog.Error("invalid jwt token")
		return nil, err
	}
	return claims, nil
}

func GetIdFromToken(r *http.Request, cfg *config.Config) (string, int) {
	var softToken string
	token := r.Header.Get("Authorization")

	if token == "" {
		return "unauthorized", http.StatusUnauthorized
	} else if strings.Contains(token, "Bearer") {
		softToken = strings.TrimPrefix(token, "Bearer ")
	} else {
		softToken = token
	}

	claims, err := ExtractClaim(softToken)
	if err != nil {
		return "unauthorized", http.StatusUnauthorized
	}

	return cast.ToString(claims["username"]), 0
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
