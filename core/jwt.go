package core

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"video-api/constants"
)

var jwtSecret = []byte("your-secret-key") // should use environment variable, for new keeping it here

func GenerateJWTForVideosSharing(videoURL string) (string, error) {
	// Set token claims
	claims := jwt.MapClaims{
		"video_url": videoURL,
		"exp":       time.Now().Add(constants.SharedLinkExpiry).Unix(),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (string, error) {
	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		videoURL, ok := claims["video_url"].(string)
		if !ok {
			return "", fmt.Errorf("invalid token claims")
		}
		return videoURL, nil
	}

	return "", fmt.Errorf("invalid token")
}
