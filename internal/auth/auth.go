package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type grantTokenResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

func FetchGrantToken(n8nUri, authToken string) (string, error) {
	url := fmt.Sprintf("http://%s/runners/auth", n8nUri)

	payload := map[string]string{"token": authToken}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch grant token: status code %d", resp.StatusCode)
	}

	var tokenResp grantTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.Data.Token, nil
}
