package cat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CatAPIValidator struct {
	client *http.Client
}

func NewCatAPIValidator() *CatAPIValidator {
	return &CatAPIValidator{
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (v *CatAPIValidator) ValidateBreed(ctx context.Context, breed string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.thecatapi.com/v1/breeds", nil)
	if err != nil {
		return false, err
	}

	resp, err := v.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var breeds []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return false, err
	}

	for _, b := range breeds {
		if b.Name == breed {
			return true, nil
		}
	}

	return false, nil
}
