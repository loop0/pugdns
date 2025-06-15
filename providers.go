package main

import (
	"log/slog"
	"os"

	"github.com/loop0/pugdns/providers/cloudflare"
	"github.com/loop0/pugdns/providers/desec"
	"github.com/loop0/pugdns/providers/ipify"
	"github.com/loop0/pugdns/providers/myipio"
	"github.com/loop0/pugdns/providers/viaip"
)

type IPAddressService interface {
	GetPublicIP() (string, error)
}

type DomainService interface {
	UpdateDomain(ip string) error
}

func getIPAddressProvider() IPAddressService {
	provider := os.Getenv("PUGDNS_IP_PROVIDER")
	if provider == "" {
		provider = "ipify"
	}
	slog.Info("Using ip address", "provider", provider)
	switch provider {
	case "ipify":
		return ipify.NewClient()
	case "myipio":
		return myipio.NewClient()
	case "viaip":
		return viaip.NewClient()
	default:
		slog.Error("Unsupported ip address provider", "provider", provider)
		os.Exit(1)
	}
	return nil
}

func getDomainProvider() DomainService {
	provider := os.Getenv("PUGDNS_DNS_PROVIDER")
	if provider == "" {
		provider = "cloudflare"
	}
	slog.Info("Using dns", "provider", provider)
	switch provider {
	case "cloudflare":
		return cloudflare.NewClient()
	case "desec":
		return desec.NewDeSECClient()
	default:
		slog.Error("Unsupported dns provider", "provider", provider)
		os.Exit(1)
	}
	return nil
}
