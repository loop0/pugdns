package viaip

import (
	"encoding/json"
	"io"
	"net/http"
)

type ViaIPClient struct {
	BaseURL string
}

type PublicIP struct {
	IP string `json:"ip"`
}

func (client *ViaIPClient) GetPublicIP() (string, error) {
	req, err := http.NewRequest("GET", client.BaseURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	publicIP := PublicIP{}
	json.Unmarshal(body, &publicIP)
	return publicIP.IP, nil
}

func NewClient() *ViaIPClient {
	return &ViaIPClient{
		"https://viaip.com.br",
	}
}
