package utils

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/wkosaibaty/lightweight-netflix/config"
	"github.com/wkosaibaty/lightweight-netflix/models"
)

func CreateJWT(user *models.User) (string, error) {
	config, _ := config.LoadConfig(".")

	privateKey, err := base64.StdEncoding.DecodeString(config.AccessTokenPrivateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = user.ID
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(config.AccessTokenExpiresIn).Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}
