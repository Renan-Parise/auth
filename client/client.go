package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Renan-Parise/auth/utils"
)

type FinancesService interface {
	CreateDefaultCategories(userID int64) error
}

type financesService struct {
	baseURL string
	client  *http.Client
}

func NewFinancesService() FinancesService {
	baseURL := os.Getenv("FINANCES_SERVICE_URL")
	return &financesService{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (fs *financesService) CreateDefaultCategories(userID int64) error {
	url := fmt.Sprintf("%s/categories/default", fs.baseURL)

	payload := map[string]int64{
		"userId": userID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	token, err := utils.GenerateServiceToken()
	if err != nil {
		return fmt.Errorf("failed to generate service token: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := fs.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call finances service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create default categories: status %d", resp.StatusCode)
	}

	return nil
}
