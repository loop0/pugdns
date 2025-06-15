package desec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/loop0/pugdns/utils"
)

type DeSECClient struct {
	APIToken string
	BaseURL  string
	Zone     string
	Domain   string
	Client   http.Client
}

type DNSRecord struct {
	Records []string `json:"records"`
}

type DNSRecordUpdate struct {
	DNSRecord
}

type DNSRecordResponse struct {
	DNSRecord
}

func NewDeSECClient() *DeSECClient {
	apiToken := utils.GetEnvOrExit("PUGDNS_DESEC_TOKEN")
	zone := utils.GetEnvOrExit("PUGDNS_ZONE")
	domain := utils.GetEnvOrExit("PUGDNS_DOMAIN")

	return &DeSECClient{
		APIToken: apiToken,
		BaseURL:  "https://desec.io/api/v1",
		Zone:     zone,
		Domain:   domain,
		Client:   *http.DefaultClient,
	}
}

func (client *DeSECClient) request(method string, url string, data interface{}, params map[string]string, reqData interface{}) error {
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

func (client *DeSECClient) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Token "+client.APIToken)
	req.Header.Set("Content-Type", "application/json")
	return client.Client.Do(req)
}

func (client *DeSECClient) UpdateDomain(ip string) error {
	name := strings.TrimSuffix(client.Domain, fmt.Sprintf(".%s", client.Zone))
	url := fmt.Sprintf("%s/domains/%s/rrsets/%s/%s/", client.BaseURL, client.Zone, name, "A")

	dnsRecordResponse := DNSRecordResponse{}
	err := client.request("GET", url, &dnsRecordResponse, nil, nil)
	if err != nil {
		return fmt.Errorf("unable to get dns record: %v", err)
	}

	// We only update the record if the ip changed
	if len(dnsRecordResponse.DNSRecord.Records) > 0 && dnsRecordResponse.DNSRecord.Records[0] != ip {
		payload := DNSRecordUpdate{
			DNSRecord: DNSRecord{
				Records: []string{ip},
			},
		}

		err := client.request("PATCH", url, nil, nil, payload)
		if err != nil {
			return fmt.Errorf("unable to update dns record: %v", err)
		}
	}

	return nil
}
