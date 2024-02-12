package main

import (
	"log/slog"
	"os"

	"github.com/loop0/pugdns/cloudflare"
	"github.com/loop0/pugdns/ipify"
)

type Config struct {
	Zone     string
	Domain   string
	ApiToken string
}

func LoadConfig() Config {
	return Config{
		os.Getenv("PUGDNS_ZONE"),
		os.Getenv("PUGDNS_DOMAIN"),
		os.Getenv("PUGDNS_CLOUDFLARE_TOKEN"),
	}
}

func main() {
	config := LoadConfig()

	ipify := ipify.NewClient()

	slog.Info("Obtaining public ip")
	ip, err := ipify.GetPublicIP()
	if err != nil {
		slog.Error("Unable to obtain public ip", "error", err)
		os.Exit(1)
	}
	slog.Info("Public", "ip", ip.IP)

	slog.Info("Updating", "domain", config.Domain)
	cloudflare := cloudflare.NewClient(config.ApiToken)
	zone, err := cloudflare.GetZoneByName(config.Zone)
	if err != nil {
		slog.Error("Unable to obtain dns zone info", "error", err)
		os.Exit(1)
	}

	dns, err := cloudflare.GetDNSRecordByName(zone, config.Domain)
	if err != nil {
		slog.Error("Unable to obtain dns record", "error", err)
		os.Exit(1)
	}

	if dns.Content != ip.IP {
		_, err = cloudflare.UpdateDNSRecord(zone, dns, config.Domain, ip.IP)
		if err != nil {
			slog.Error("Unable to update dns record", "error", err)
			os.Exit(1)
		}
		slog.Info("Updated", "domain", config.Domain, "ip", ip.IP)
	} else {
		slog.Info("No changes to ip address, no update required")
	}
}
