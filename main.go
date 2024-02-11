package main

import (
	"log"
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

	log.Println("Obtaining public ip")
	ip, err := ipify.GetPublicIP()
	if err != nil {
		log.Fatalf("Unable to obtain public ip: %v", err)
	}
	log.Printf("Public ip is %v\n", ip.IP)

	log.Printf("Updating domain %v\n", config.Domain)
	cloudflare := cloudflare.NewClient(config.ApiToken)
	zone, err := cloudflare.GetZoneByName(config.Zone)
	if err != nil {
		log.Fatalf("Unable to obtain dns zone info: %v", err)
	}

	dns, err := cloudflare.GetDNSRecordByName(zone, config.Domain)
	if err != nil {
		log.Fatalf("Unable to obtain dns record: %v", err)
	}

	if dns.Content != ip.IP {
		_, err = cloudflare.UpdateDNSRecord(zone, dns, config.Domain, ip.IP)
		if err != nil {
			log.Fatalf("Unable to update dns record: %v", err)
		}
		log.Printf("Domain %v updated with ip %v", config.Domain, ip.IP)
	} else {
		log.Printf("No changes to ip address, no update required")
	}
}
