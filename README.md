# pugdns

![Yoda](images/yoda.png)

pugdns updates a domain on the dns provider of choice with your public ip address.

## Why

I developed this small application to relearn Go and because I needed
something to update my home vpn domain with my public address when it
changes.

## How to use

pugdns gets the required configuration from environment variables as follows:

```sh
PUGDNS_ZONE=example.com  # The zone name where the domain will be updated
PUGDNS_DOMAIN=vpn.example.com  # The domain name to be updated with your public ip
PUGDNS_DNS_PROVIDER=desec  # The dns provider to be used when updating the dns record
PUGDNS_DESEC_TOKEN=<API_TOKEN>  # The desc api token with read/write access to the zone
PUGDNS_IP_PROVIDER=ipify  # The ip provider to resolve your public ip
```

Once you have those env vars set you can just run:

```sh
$ pugdns
2024/02/11 16:37:57 INFO Obtaining public ip
2024/02/11 16:37:57 INFO Public ip=xxx.xxx.xxx.xxx
2024/02/11 16:37:57 INFO Updating domain=vpn.example.com
2024/02/11 16:37:57 INFO Updated domain=vpn.example.com ip=xxx.xxx.xxx.xxx
```

## Providers

Some providers will require environment variables to be set

### IP Address

IP Address provider is selected via the `PUGDNS_IP_PROVIDER` environment variable. Default is `ipify`

This is the list of currently supported providers for public ip address resolution:

| Provider | Website | Required Env Vars |
| - | - | - |
| `ipify` | https://www.ipify.org | None |
| `myipio` | https://www.my-ip.io | None |
| `viaip` | https://viaip.com.br | None |

### DNS

IP Address provider is selected via the `PUGDNS_DNS_PROVIDER` environment variable. Default is `cloudflare`
Note that current providers only support updating subdomains, updating the apex is currently not supported.
The application also assumes that you already have your dns record created at your provider.

This is the list of currently supported providers for DNS:

| Provider | Website | Required Env Vars |
| - | - | - |
| `desec` | https://desec.io | `PUGDNS_DESEC_TOKEN` |
| `cloudflare` | https://www.cloudflare.com | `PUGDNS_CLOUDFLARE_TOKEN` |


### How to add providers

There are currently only 2 types of providers required by pugdns to work which are defined by the following interfaces at `providers.go`:

```go
type IPAddressService interface {
	GetPublicIP() (string, error)
}

type DomainService interface {
	UpdateDomain(ip string) error
}
```

To add a new provider you just need to implement its interface and add it to `providers.go` function's `getIPAddressProvider` or `getDomainProvider`. For existing provider's implementation look into the `providers` folder.

## Usage examples
### kubernetes

If you want to run pugdns in your kubernetes cluster take a look into `examples/k8s-cronjob.yaml` for a manifest for a cronjob that runs every hour to keep your domain update with you public ip address.
