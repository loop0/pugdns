package ipify

import (
	"encoding/json"
	"io"
	"net/http"
)

type IPifyClient struct {
	BaseURL string
}

type PublicIP struct {
	IP string `json:"ip"`
}

func (client *IPifyClient) GetPublicIP() (PublicIP, error) {
	resp, err := http.Get(client.BaseURL)
	if err != nil {
		return PublicIP{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PublicIP{}, err
	}

	publicIP := PublicIP{}
	json.Unmarshal(body, &publicIP)
	return publicIP, nil
}

func NewClient() *IPifyClient {
	return &IPifyClient{
		"https://api.ipify.org?format=json",
	}
}
