# Check Expiring Certs

Check Expiring Certs is a fast and simple tool designed for efficiently
monitoring HTTPS certificates. It warns if an SSL/TLS certificate on a
specified hosts is due to expire within a configurable number of days.

## Installation

1. **Using Go:**
   - Install Go (version 1.x or later recommended).
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
check-expiring-certs [options] <host:port> [<host:port> ...]
```

It exits with return code 1 if any certificates are expiring soon or
connections failed.

The tool accepts `hostname:port` arguments, defaulting to port 443 if omitted,
and allows `server_name:hostname:port` to specify a different server name for
TLS SNI.

### Options

- `-warn N`: Warn if a certificate is expiring within N days (default: 7).
- `-timeout DURATION`: Set a timeout for each connection (default: 2s).
- `-concurrency N`: Specify the number of concurrent checks (default: 128).

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
