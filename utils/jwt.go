package utils

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/wkosaibaty/lightweight-netflix/config"
	"github.com/wkosaibaty/lightweight-netflix/models"
)

func CreateJWT(user *models.User) (string, error) {
	config, _ := config.LoadConfig(".")

	privateKey, err := base64.StdEncoding.DecodeString(config.AccessTokenPrivateKey)
	if err != nil {
		return "", err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = user.ID
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(config.AccessTokenExpiresIn).Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateToken(token string) (interface{}, error) {
	config, _ := config.LoadConfig(".")

	publicKey, err := base64.StdEncoding.DecodeString(config.AccessTokenPublicKey)
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return "", err
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("Invalid access token")
		}
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("Invalid access token")
	}

	return claims["sub"], nil
}
