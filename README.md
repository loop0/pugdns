# pugdns

![Yoda](images/yoda.png)

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
2024/02/11 16:37:57 INFO Obtaining public ip
2024/02/11 16:37:57 INFO Public ip=xxx.xxx.xxx.xxx
2024/02/11 16:37:57 INFO Updating domain=vpn.example.com
2024/02/11 16:37:57 INFO Updated domain=vpn.example.com ip=xxx.xxx.xxx.xxx
```

## kubernetes

If you want to run pugdns in your kubernetes cluster take a look into `examples/k8s-cronjob.yaml` for a manifest for a cronjob that runs every hour to keep your domain update with you public ip address.