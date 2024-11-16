package utils

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/Renan-Parise/auth/entities"
	"github.com/Renan-Parise/auth/errors"
)

func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func GenerateCode(length int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(code)
}

func SendEmail(email entities.Email) error {
	mailServiceURL := GetMailServiceURL() + "/mail/send"

	jsonData, err := json.Marshal(email)
	if err != nil {
		return errors.NewServiceError("Failed to send email: " + err.Error())
	}

	req, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.NewServiceError("Failed to send email: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		return errors.NewServiceError("Failed to send email: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.NewServiceError("Failed to send email, status code: " + resp.Status)
	}

	return nil
}
