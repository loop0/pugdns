package myipio

import (
	"encoding/json"
	"io"
	"net/http"
)

type MyIPIOClient struct {
	BaseURL string
}

type PublicIP struct {
	IP string `json:"ip"`
}

func (client *MyIPIOClient) GetPublicIP() (string, error) {
	resp, err := http.Get(client.BaseURL)
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

func NewClient() *MyIPIOClient {
	return &MyIPIOClient{
		"https://api.my-ip.io/v2/ip.json",
	}
}
