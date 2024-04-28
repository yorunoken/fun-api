package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
)

func RandomImage(path string) (*string, error) {
	imageFiles, err := os.ReadDir("./media/" + path)

	if err != nil {
		return nil, errors.New("error reading directory")
	}

	if len(imageFiles) == 0 {
		return nil, errors.New("no images found in directory")
	}

	randomIndex := rand.Intn(len(imageFiles))
	randomImage := imageFiles[randomIndex].Name()

	imagePath := fmt.Sprintf("/media/%s/%s", path, randomImage)

	return &imagePath, nil
}
