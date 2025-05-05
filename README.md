[![License](https://img.shields.io/github/license/ei-grad/check-expiring-certs)](LICENSE)
[![Workflow Status](https://github.com/ei-grad/check-expiring-certs/actions/workflows/release.yml/badge.svg)](https://github.com/ei-grad/check-expiring-certs/actions/workflows/release.yml)
[![Go Coverage](https://github.com/ei-grad/check-expiring-certs/wiki/coverage.svg)](https://raw.githack.com/wiki/ei-grad/check-expiring-certs/coverage.html)
[![Go Report](https://goreportcard.com/badge/github.com/ei-grad/check-expiring-certs)](https://goreportcard.com/report/github.com/ei-grad/check-expiring-certs)
[![Latest Release](https://img.shields.io/github/v/release/ei-grad/check-expiring-certs)](https://github.com/ei-grad/check-expiring-certs/releases/latest)
[![Downloads](https://img.shields.io/github/downloads/ei-grad/check-expiring-certs/total)](https://github.com/ei-grad/check-expiring-certs/graphs/traffic)
[![Contributors](https://img.shields.io/github/contributors/ei-grad/check-expiring-certs)](https://github.com/ei-grad/check-expiring-certs/graphs/contributors)
[![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/eigrad)](https://x.com/eigrad)

![GPT-4o refined catchy banner](https://repository-images.githubusercontent.com/93572949/43a8be00-1330-410e-91e2-640c898c002b)

# Check Expiring Certs

Check Expiring Certs is a fast and simple tool designed for efficiently
monitoring HTTPS certificates expiration. It warns if an SSL/TLS certificate on
a specified hosts is due to expire within a configurable number of days.

## Installation

1. **Using Go:**
   - Install Go, version 1.21+ required.
   - Run the following command:
     ```
     go install github.com/ei-grad/check-expiring-certs@latest
     ```

2. **Alternative Method - Github Releases:**
   1. **Download the archive:**
      - Visit the [releases page](https://github.com/ei-grad/check-expiring-certs/releases/latest).
      - Choose and download the file for your operating system and platform
        (amd64 or arm64). For Mac/Linux, it's a `.tar.gz` file, and for
        Windows, it's a `.zip` file.
   2. **Extract the Archive:**
      - For Mac/Linux: Use the command `tar -xzf [filename].tar.gz` in the terminal.
      - For Windows: Right-click the `.zip` file and select "Extract All...".
   3. **Using the Executable:**
      - For Mac/Linux: In the terminal, navigate to the extracted folder and run
        `./check-expiring-certs`.
      - For Windows: Open Command Prompt, navigate to the extracted folder, and
        run `check-expiring-certs.exe`.

## Usage

Run the command with the following syntax:
```
check-expiring-certs [options] <host:port>...
```

It exits with return code 1 if any certificates are expiring soon or
connections failed, and writes a message like 'Certificate for `host`
(`common-name`) expires in `time-until-expiration`'. Connection errors are also
logged to stdout.

The tool accepts `host:port` arguments, defaulting to port 443 if omitted, and
allows `server_name:host:port` to specify a different server name for TLS SNI.

### Options

- `-warn N`: Warn if a certificate is expiring within N days (default: 7).
- `-timeout DURATION`: Set a timeout for each connection (default: 2s).
- `-c N`: Specify the number of concurrent checks (default: 128).

### Examples

1. Checking multiple domains:
   ```
   check-expiring-certs -warn 14 google.com github.com
   ```

2. Custom timeout and checking internal IPs:
   ```
   check-expiring-certs -warn 30 -timeout 5s 10.0.0.1:8443 192.168.0.5:9443
   ```

3. Checking an IPv6 address:
   ```
   check-expiring-certs [2a02:6b8:a::a]:443
   ```

4. Custom IP address to connect:
   ```
   check-expiring-certs github.com:1.2.3.4:443
   ```

## License

This project is licensed under the [MIT License](LICENSE).
