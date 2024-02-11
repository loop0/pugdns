# pugdns

pugdns updates a domain on cloudflare with your public ip address.

## Why

I developed this small application to relearn Go and because I needed
something to update my home vpn domain with my public address when it
changes.

## How to use

pugdns gets the required configuration from environment variables as follows:

```sh
PUGDNS_ZONE=example.com  # The zone name where the domain will be updated
PUGDNS_DOMAIN=vpn.example.com  # The domain name to be updated with your public ip
PUGDNS_CLOUDFLARE_TOKEN=<API_TOKEN>  # The cloudflare api token with read/write access to the zone
```

Once you have those env vars set you can just run:

```sh
$ pugdns
2024/02/11 16:37:57 Obtaining public ip
2024/02/11 16:37:57 Public ip is xxx.xxx.xxx.xxx
2024/02/11 16:37:57 Updating domain vpn.example.com
2024/02/11 16:37:57 Domain vpn.example.com updated with ip xxx.xxx.xxx.xxx
```

## TODO
- [ ] Dockerfile to create container
- [ ] Instructions on how to run on kubernetes