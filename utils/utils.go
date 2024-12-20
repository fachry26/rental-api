package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"rental-api/models"
)

func RespondError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"status":  "error",
		"message": message,
	})
}

func RespondJSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status": "success",
		"data":   data,
	})
}

func GenerateJWT(user *models.User) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(secretKey) == 0 {
		log.Fatal("JWT_SECRET_KEY not set in environment variables")
		return "", fmt.Errorf("secret key not set")
	}

	expirationDuration, err := time.ParseDuration(os.Getenv("JWT_EXPIRATION_DURATION"))
	if err != nil {
		expirationDuration = 24 * time.Hour
	}

	claims := &jwt.StandardClaims{
		Subject:   fmt.Sprintf("%d", user.ID),
		ExpiresAt: time.Now().Add(expirationDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

func ParseJWT(tokenStr string) (*jwt.StandardClaims, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(secretKey) == 0 {
		log.Fatal("JWT_SECRET_KEY not set in environment variables")
		return nil, fmt.Errorf("secret key not set")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid or expired token: %v", err)
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}

func ExtractUserIDFromJWT(c *gin.Context) (uint, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return 0, fmt.Errorf("token is missing in the Authorization header")
	}

	claims, err := ParseJWT(tokenString)
	if err != nil {
		return 0, fmt.Errorf("failed to parse JWT: %v", err)
	}

	userIDStr := claims.Subject
	if userIDStr == "" {
		return 0, fmt.Errorf("invalid user ID in token")
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to convert user ID to uint: %v", err)
	}

	return uint(userID), nil
}
