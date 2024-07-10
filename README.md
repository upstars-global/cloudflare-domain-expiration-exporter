# cloudflare-domain-expiration-exporter

Domain expiration exporter for Cloudflare.

## Usage

```bash
$ CF_API_KEY="key1,key2" cloudflare-domain-expiration-exporter
```

## Port

API server listens `8080` port.

## Environment Variables

- `CF_API_KEYS`: Cloudflare API keys. Multiple keys can be specified by separating them with commas.

## Metrics

- `domain_expiration_checker_result`: Number of days until the domain expires.
  - contains two labels: `domain` and `status`.
    - `domain`: Domain name.
    - `status`: Status of the domain expiration.
      - `ok`: The domain is not expired.
      - `unknown`: The domain expiration status is unknown.

## Set manual expiration

You can set the expiration date manually by setting domain expirations in file `manual_expiration.yaml`.
To do that you can simply copy the `manual_expiration.example.yaml` file and set the expiration date for the domains.

## License

MIT