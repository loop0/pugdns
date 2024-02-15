package cloudflare

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/loop0/pugdns/utils"
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
	Zone      string
	Domain    string
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

	resp, err := client.do(req)
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

func (client *CloudflareClient) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", client.APIToken))
	req.Header.Set("Content-Type", "application/json")
	return client.Client.Do(req)
}

func (client *CloudflareClient) getZoneByName(name string) (Zone, error) {
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

func (client *CloudflareClient) getDNSRecordByName(zone Zone, name string) (DNSRecord, error) {
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

func (client *CloudflareClient) updateDNSRecord(zone Zone, record DNSRecord, name string, content string) (DNSRecord, error) {
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

func (client *CloudflareClient) UpdateDomain(ip string) error {
	zone, err := client.getZoneByName(client.Zone)
	if err != nil {
		return fmt.Errorf("unable to obtain dns zone info: %v", err)
	}

	dns, err := client.getDNSRecordByName(zone, client.Domain)
	if err != nil {
		return fmt.Errorf("unable to obtain dns record: %v", err)
	}

	if dns.Content != ip {
		_, err = client.updateDNSRecord(zone, dns, client.Domain, ip)
		if err != nil {
			return fmt.Errorf("unable to update dns record: %v", err)
		}
	}

	return nil
}

func NewClient() *CloudflareClient {
	apiToken := utils.GetEnvOrExit("PUGDNS_CLOUDFLARE_TOKEN")
	zone := utils.GetEnvOrExit("PUGDNS_ZONE")
	domain := utils.GetEnvOrExit("PUGDNS_DOMAIN")

	return &CloudflareClient{
		*http.DefaultClient,
		apiToken,
		"https://api.cloudflare.com",
		"client/v4",
		zone,
		domain,
	}
}
