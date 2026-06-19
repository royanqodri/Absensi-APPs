package util

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Base64ToImage converts a base64 string to an image file
// func Base64ToImage(base64String string) (string, string, error) {
// 	if !ValidateBase64(base64String) {
// 		return "", "", fmt.Errorf("invalid base64 string for image")
// 	}

// 	// imageType, decodedData, err := decodeBase64(base64String)
// 	// if err != nil {
// 	// 	return "", "", fmt.Errorf("error converting base64 to image: %v", err)
// 	// }

// 	// img, format, err := image.Decode(bytes.NewReader(decodedData))
// 	// if err != nil {
// 	// 	return "", "", fmt.Errorf("error decoding image: %v", err)
// 	// }

// 	// randomFileName := generateRandomFileName() + "." + imageType
// 	// fileLocalPath := config.Get().FileLocalPath
// 	// fileLocal := filepath.Join(fileLocalPath, randomFileName)

// 	// outputFile, err := createFile(fileLocal)
// 	// if err != nil {
// 	// 	return "", "", fmt.Errorf("error creating output file: %v", err)
// 	// }
// 	// defer outputFile.Close()

// 	// if format == "jpeg" {
// 	// 	err = jpeg.Encode(outputFile, img, nil)
// 	// } else if format == "png" {
// 	// 	err = png.Encode(outputFile, img)
// 	// }
// 	// if err != nil {
// 	// 	return "", "", fmt.Errorf("error encoding image: %v", err)
// 	// }

// 	// return randomFileName, "", nil
// }

// Base64ToSound converts a base64 string to a sound file
// func Base64ToSound(base64String string) (string, string, error) {
// 	if !ValidateBase64(base64String) {
// 		return "", "", fmt.Errorf("invalid base64 string for sound")
// 	}

// 	soundType, decodedData, err := decodeBase64(base64String)
// 	if err != nil {
// 		return "", "", fmt.Errorf("error converting base64 to sound: %v", err)
// 	}

// 	fileName := GenerateUUIDFileName() + "." + soundType
// 	fileLocalPath := config.Get().FileLocalPath
// 	fileLocal := filepath.Join(fileLocalPath, fileName)

// 	err = os.WriteFile(fileLocal, decodedData, 0644)
// 	if err != nil {
// 		return "", "", fmt.Errorf("error writing sound file: %v", err)
// 	}

// 	return fileName, fileLocal, nil
// }

// decodeBase64 decodes a base64 string and determines its file type
func decodeBase64(base64String string) (string, []byte, error) {
	var fileType string
	var decodedData []byte

	switch {
	case strings.HasPrefix(base64String, "data:image/png;base64,"):
		fileType = "png"
		base64String = strings.TrimPrefix(base64String, "data:image/png;base64,")
	case strings.HasPrefix(base64String, "data:image/jpg;base64,"):
		fallthrough
	case strings.HasPrefix(base64String, "data:image/jpeg;base64,"):
		fileType = "jpeg"
		base64String = strings.TrimPrefix(base64String, "data:image/jpeg;base64,")
	case strings.HasPrefix(base64String, "data:audio/mp3;base64,"):
		fileType = "mp3"
		base64String = strings.TrimPrefix(base64String, "data:audio/mp3;base64,")
	case strings.HasPrefix(base64String, "data:audio/mpeg;base64,"):
		fileType = "mp3"
		base64String = strings.TrimPrefix(base64String, "data:audio/mpeg;base64,")
	default:
		return "", nil, fmt.Errorf("unsupported file format")
	}

	decodedData, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", nil, fmt.Errorf("error decoding base64 string: %v", err)
	}

	return fileType, decodedData, nil
}

// generateRandomFileName generates a random file name
func generateRandomFileName() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// generateUUIDFileName generates a UUID based file name
func GenerateUUIDFileName() string {
	u := uuid.New()
	return u.String()
}

// createFile creates a new file
func createFile(filePath string) (*os.File, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %v", err)
	}
	return file, nil
}

// DecodeBase64ToFile decodes a base64 string and writes it to a file
func DecodeBase64ToFile(base64String, filePath string) error {
	// Decode base64 string
	decoded, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return fmt.Errorf("error decoding base64 string: %v", err)
	}

	// Write decoded data to file
	err = ioutil.WriteFile(filePath, decoded, 0644)
	if err != nil {
		return fmt.Errorf("error writing decoded data to file: %v", err)
	}

	return nil
}

// ValidateBase64 checks if the input string is a valid base64 encoded string.
func ValidateBase64(input string) bool {
	// Regular expression to match base64 string
	base64Regex := regexp.MustCompile(`^data:(image\/(png|jpeg)|audio\/(mp3|mpeg));base64,([a-zA-Z0-9+/=]+)$`)

	// Check if input matches base64 regex
	return base64Regex.MatchString(input)
}

// Base64Compare compares two Base64-encoded strings and returns true if they are identical after decoding.
func Base64Compare(encoded1, encoded2 string) (bool, error) {
	// Decode the first Base64 string
	data1, err := base64.StdEncoding.DecodeString(encoded1)
	if err != nil {
		return false, errors.New("invalid Base64 string in the first argument")
	}

	// Decode the second Base64 string
	data2, err := base64.StdEncoding.DecodeString(encoded2)
	if err != nil {
		return false, errors.New("invalid Base64 string in the second argument")
	}

	// Compare the decoded byte slices
	return bytes.Equal(data1, data2), nil
}
