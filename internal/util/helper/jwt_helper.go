package helper

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID        string `json:"user_id"`
	Email         string `json:"email"`
	RoleId        string `json:"role_id"`
	DepartementId string `json:"department_id"`
	IsActive      bool   `json:"is_active"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, email, roleID, departementId string, isActive bool) (string, int64, error) {
	// Load environment variables
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	jwtIssuer := os.Getenv("JWT_ISSUER")
	expHoursStr := os.Getenv("JWT_EXPIRATION_HOURS")

	jwtExpHours, err := strconv.Atoi(expHoursStr)
	if err != nil || jwtExpHours <= 0 {
		jwtExpHours = 1
		log.Println("⚠️ JWT_EXPIRATION_HOURS tidak ditemukan, default ke 1 jam")
	}

	expiresIn := int64(jwtExpHours * 3600)

	claims := &JWTClaims{
		UserID:        userID,
		Email:         email,
		RoleId:        roleID,
		DepartementId: departementId,
		IsActive:      isActive,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtExpHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    jwtIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", 0, err
	}

	return signedToken, expiresIn, nil
}
