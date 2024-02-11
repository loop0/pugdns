package cloudflare

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Zone struct {
	ID string `json:"id"`
}

type ZoneResponse struct {
	Zones []Zone `json:"result"`
}

type DNSRecord struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type DNSRecordUpdate struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

type DNSRecordsResponse struct {
	DNSRecords []DNSRecord `json:"result"`
}

type DNSRecordResponse struct {
	DNSRecord DNSRecord `json:"result"`
}

type CloudflareClient struct {
	Client    http.Client
	APIToken  string
	BaseURL   string
	APIPrefix string
}

func (client *CloudflareClient) request(method string, url string, data interface{}, params map[string]string, reqData interface{}) error {
	var body io.Reader
	if reqData != nil {
		reqBody, err := json.Marshal(reqData)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(reqBody)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	query := req.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error calling cloudflare api: %v", string(respBody))
	}

	json.Unmarshal(respBody, data)
	return nil
}

func (client *CloudflareClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", client.APIToken))
	req.Header.Set("Content-Type", "application/json")
	return client.Client.Do(req)
}

func (client *CloudflareClient) GetZoneByName(name string) (Zone, error) {
	url := fmt.Sprintf("%s/%v/zones", client.BaseURL, client.APIPrefix)
	params := map[string]string{"name": name}
	data := ZoneResponse{}

	err := client.request("GET", url, &data, params, nil)
	if err != nil {
		return Zone{}, err
	}

	if len(data.Zones) < 1 {
		return Zone{}, errors.New("missing zone info")
	}

	return data.Zones[0], nil
}

func (client *CloudflareClient) GetDNSRecordByName(zone Zone, name string) (DNSRecord, error) {
	url := fmt.Sprintf("%v/%v/zones/%v/dns_records", client.BaseURL, client.APIPrefix, zone.ID)
	params := map[string]string{"name": name}
	data := DNSRecordsResponse{}

	err := client.request("GET", url, &data, params, nil)
	if err != nil {
		return DNSRecord{}, err
	}

	if len(data.DNSRecords) < 1 {
		return DNSRecord{}, errors.New("missing dns record")
	}
	return data.DNSRecords[0], nil
}

func (client *CloudflareClient) UpdateDNSRecord(zone Zone, record DNSRecord, name string, content string) (DNSRecord, error) {
	url := fmt.Sprintf("%v/%v/zones/%v/dns_records/%v", client.BaseURL, client.APIPrefix, zone.ID, record.ID)
	data := DNSRecordResponse{}
	reqData := DNSRecordUpdate{
		name,
		content,
		"A",
	}

	err := client.request("PATCH", url, &data, nil, reqData)
	if err != nil {
		return DNSRecord{}, err
	}

	return data.DNSRecord, nil
}

func NewClient(apiToken string) *CloudflareClient {
	return &CloudflareClient{
		*http.DefaultClient,
		apiToken,
		"https://api.cloudflare.com",
		"client/v4",
	}
}
