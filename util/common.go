package util

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
)

var Validate = validator.New()

func GoDotEnvVariable(key string, fallBack string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	// return os.Getenv(key)
	return fallBack
}

func AtoI(s string, v int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return v
	}
	return i
}

func AtoF(s string, v float64) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return v
	}
	return f
}

func CalcFileSizeFromString(data string) int {
	d := len(data)
	// count how many trailing '=' there are (if any)
	eq := 0
	if d >= 2 {
		if data[d-1] == '=' {
			eq++
		}
		if data[d-2] == '=' {
			eq++
		}
		d -= eq
	}
	return (d*3 - eq) / 4
}

// Check images base64
func CheckImage(image string, mimeType string) (string, error) {
	_, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		return "", errors.New(fmt.Sprint("broken image: ", err.Error()))
	}
	rawImage := fmt.Sprintf("%s%s", image, mimeType)
	if CalcFileSizeFromString(rawImage) > 2000000 { // 2MB
		return "", errors.New("image over size")
	}
	return image, nil
}

// Check Task Status
func CheckTaskStatus(status string) bool {
	if status == "IN_PROGRESS" || status == "COMPLETED" {
		return true
	}
	return false
}

func ValidateStruct(data interface{}) string {
	var result string
	err := Validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			result = strings.ToLower(err.StructField())
		}
	}
	return result
}
