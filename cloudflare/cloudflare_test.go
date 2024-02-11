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
	}
}

func TestGetZoneByName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		expected := "/client/v4/zones?name=loop0.sh"
		if req.URL.String() != expected {
			t.Fatalf("Request path mismatch: %v != %v", expected, req.URL.String())
		}
		rw.Write([]byte(`{"result": [{"id": "test"}]}`))
	}))

	client := NewTestClient(server.URL)

	zone, _ := client.GetZoneByName("loop0.sh")
	if zone.ID != "test" {
		t.Fail()
	}

}

func TestGetDNSRecordByName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		expected := "/client/v4/zones/test/dns_records?name=vpn.loop0.sh"
		if req.URL.String() != expected {
			t.Fatalf("Request path mismatch: %v != %v", expected, req.URL.String())
		}
		rw.Write([]byte(`{"result": [{"id": "test"}]}`))
	}))

	client := NewTestClient(server.URL)
	zone := Zone{"test"}
	dns, _ := client.GetDNSRecordByName(zone, "vpn.loop0.sh")
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
	updated, _ := client.UpdateDNSRecord(zone, dns, "vpn.loop0.sh", "127.0.0.1")
	if updated.ID != "test" {
		t.Fail()
	}
}
