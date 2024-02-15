package main

import (
	"os"

	"github.com/loop0/pugdns/providers/cloudflare"
	"github.com/loop0/pugdns/providers/ipify"
)

type IPAddressService interface {
	GetPublicIP() (string, error)
}

type DomainService interface {
	UpdateDomain(ip string) error
}

func getIPAddressProvider() IPAddressService {
	provider := os.Getenv("PUGDNS_IP_PROVIDER")
	switch provider {
	default:
		return ipify.NewClient()
	}
}

func getDomainProvider() DomainService {
	provider := os.Getenv("PUGDNS_DNS_PROVIDER")
	switch provider {
	default:
		return cloudflare.NewClient()
	}
}
