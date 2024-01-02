package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
)

func main() {

	exitcode := 0

	warning_period := flag.Int("warn", 7, "warning period in days")
	timeout := flag.Duration("timeout", 2*time.Second, "timeout for connection")
	concurrency := flag.Int("c", 128, "number of concurrent checks")
	flag.Parse()

	// endpoints to check
	endpoints := flag.Args()

	// semaphore to limit concurrency to a reasonable number
	semaphore := make(chan struct{}, *concurrency)

	dialer := &net.Dialer{
		Timeout: *timeout,
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(endpoints))

	warn_if_expired_at := time.Now().AddDate(0, 0, *warning_period)

	for _, i := range endpoints {

		// sleep 1ms to avoid hitting DNS resolver limits
		time.Sleep(time.Millisecond)

		// acquire semaphore
		semaphore <- struct{}{}

		go func(i string) {

			// mark as done when we're finished
			defer wg.Done()

			// release semaphore
			defer func() { <-semaphore }()

			is_expired, err := checkHost(dialer, i, warn_if_expired_at)
			if err != nil {
				fmt.Printf("can't check %s: %s\n", i, err)
				exitcode = 1
			} else if is_expired {
				exitcode = 1
			}

		}(i)
	}

	wg.Wait()

	os.Exit(exitcode)
}

var addrOverride = regexp.MustCompile(`^([^:]+):(((\[[0-9a-f:]+\])|([^:]+)):\d+)$`)

func checkHost(
	dialer *net.Dialer,
	host string,
	warn_if_expired_at time.Time,
) (is_expired bool, err error) {
	config := tls.Config{
		InsecureSkipVerify: true,
	}
	if !strings.Contains(host, ":") {
		host = host + ":443"
	} else if match := addrOverride.FindStringSubmatch(host); match != nil {
		fmt.Printf("Match %s: %s\n", host, match)
		config.ServerName = match[1]
		host = match[2]
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", host, &config)
	if err != nil {
		return
	}
	conn.Close()
	for _, cert := range conn.ConnectionState().PeerCertificates {
		if warn_if_expired_at.After(cert.NotAfter) {
			is_expired = true
			fmt.Printf("Certificate for %s (%s) expires in %s\n",
				host, cert.Subject.CommonName,
				humanize.Time(cert.NotAfter))
		}
	}
	return
}
