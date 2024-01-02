# Check Expiring Certs

This tool warns if an SSL/TLS certificate on a specified host is due to expire
within a configurable number of days.

## Usage

```
check-expiring-certs [options] <host:port> [<host:port> ...]
```

Options:

- `-warn N` - Number of days before expiration to warn about (default 7)
- `-timeout DURATION` - Timeout for each connection (default 2s)
- `-concurrency N` - Number of concurrent checks (default 128)

Examples:

```
check-expiring-certs -warn 14 google.com:443 github.com:443
check-expiring-certs -warn 30 -timeout 5s 10.0.0.1:8443 192.168.0.5:9443
check-expiring-certs [2a02:6b8:a::a]:443
```

Exits with return code 1 if any certificates are expiring soon or connections
failed.

## Installation

```
go install github.com/ei-grad/check-expiring-certs@latest
```

## License

MIT License - see [LICENSE](LICENSE) for more details.
