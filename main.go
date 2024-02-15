package main

import (
	"log/slog"
	"os"

	"github.com/loop0/pugdns/utils"
)

func main() {
	domain := utils.GetEnvOrExit("PUGDNS_DOMAIN")

	ipProvider := getIPAddressProvider()

	slog.Info("Obtaining public ip")
	ip, err := ipProvider.GetPublicIP()
	if err != nil {
		slog.Error("Unable to obtain public ip", "error", err)
		os.Exit(1)
	}
	slog.Info("Public", "ip", ip)

	slog.Info("Updating", "domain", domain)
	domainProvider := getDomainProvider()
	err = domainProvider.UpdateDomain(ip)
	if err != nil {
		slog.Error("Update to update ip address on given provider", "error", err)
		os.Exit(1)
	}
	slog.Info("Updated", "domain", domain, "ip_address", ip)
}
