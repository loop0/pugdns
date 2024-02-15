package cloudflare

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewTestClient(url string) CloudflareClient {
	return CloudflareClient{
		*http.DefaultClient,
		"token",
		url,
		"client/v4",
		"example.com",
		"vpn.example.com",
	}
}

func TestGetZoneByName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		expected := "/client/v4/zones?name=example.com"
		if req.URL.String() != expected {
			t.Fatalf("Request path mismatch: %v != %v", expected, req.URL.String())
		}
		rw.Write([]byte(`{"result": [{"id": "test"}]}`))
	}))

	client := NewTestClient(server.URL)

	zone, _ := client.getZoneByName("example.com")
	if zone.ID != "test" {
		t.Fail()
	}

}

func TestGetDNSRecordByName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		expected := "/client/v4/zones/test/dns_records?name=vpn.example.com"
		if req.URL.String() != expected {
			t.Fatalf("Request path mismatch: %v != %v", expected, req.URL.String())
		}
		rw.Write([]byte(`{"result": [{"id": "test"}]}`))
	}))

	client := NewTestClient(server.URL)
	zone := Zone{"test"}
	dns, _ := client.getDNSRecordByName(zone, "vpn.example.com")
	if dns.ID != "test" {
		t.Fail()
	}
}

func TestUpdateDNSRecord(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		expected := "/client/v4/zones/test/dns_records/test"
		if req.URL.String() != expected {
			t.Fatalf("Request path mismatch: %v != %v", expected, req.URL.String())
		}
		rw.Write([]byte(`{"result": {"id": "test"}}`))
	}))

	client := NewTestClient(server.URL)
	zone := Zone{"test"}
	dns := DNSRecord{"test", "127.0.0.1"}
	updated, _ := client.updateDNSRecord(zone, dns, "vpn.example.com", "127.0.0.1")
	if updated.ID != "test" {
		t.Fail()
	}
}
