package utils

import (
	"encoding/base64"
	"errors"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/wkosaibaty/lightweight-netflix/config"
)

func UploadImage(base64String string) (string, error) {
	config, _ := config.LoadConfig(".")
	path := config.ImagesPath

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", errors.New("Could not create image file")
	}

	bytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", errors.New("Invalid image encoding")
	}

	mimeType := http.DetectContentType(bytes)
	if !strings.HasPrefix(mimeType, "image") {
		return "", errors.New("Invalid image")
	}

	exts, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(exts) == 0 {
		return "", errors.New("Invalid image")
	}

	fileName := path + uuid.New().String() + exts[0]
	file, err := os.Create(fileName)
	if err != nil {
		return "", errors.New("Could not create image file")
	}
	defer file.Close()

	if _, err := file.Write(bytes); err != nil {
		return "", errors.New("Could not create image file")
	}
	if err := file.Sync(); err != nil {
		return "", errors.New("Could not create image file")
	}

	return fileName, nil
}
