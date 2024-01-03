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

	warning_period := flag.Int("warn", 7, "warning period in days")
	timeout := flag.Duration("timeout", 2*time.Second, "timeout for connection")
	concurrency := flag.Int("c", 128, "number of concurrent checks")
	flag.Parse()

	// endpoints to check
	endpoints := flag.Args()

	dialer := &net.Dialer{
		Timeout: *timeout,
	}

	warn_if_expired_at := time.Now().AddDate(0, 0, *warning_period)

	checker := NewSimpleHostChecker(dialer, warn_if_expired_at)

	os.Exit(RunChecks(checker, endpoints, *concurrency))
}

type HostChecker interface {
	CheckHost(host string) (bool, error)
}

func RunChecks(
	checker HostChecker,
	endpoints []string,
	concurrency int,
) (exitcode int) {

	// semaphore to limit concurrency to a reasonable number
	semaphore := make(chan struct{}, concurrency)

	wg := new(sync.WaitGroup)
	wg.Add(len(endpoints))

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

			is_expired, err := checker.CheckHost(i)
			if err != nil {
				fmt.Printf("can't check %s: %s\n", i, err)
				exitcode = 1
			} else if is_expired {
				exitcode = 1
			}

		}(i)
	}

	wg.Wait()

	return exitcode

}

var addrOverride = regexp.MustCompile(`^([^:]+):(((\[[0-9a-f:]+\])|([^:]+)):\d+)$`)

type SimpleHostChecker struct {
	dialer             *net.Dialer
	warn_if_expired_at time.Time
}

func NewSimpleHostChecker(
	dialer *net.Dialer,
	warn_if_expired_at time.Time,
) *SimpleHostChecker {
	return &SimpleHostChecker{
		dialer:             dialer,
		warn_if_expired_at: warn_if_expired_at,
	}
}

func (c *SimpleHostChecker) CheckHost(
	host string,
) (is_expired bool, err error) {

	config := tls.Config{
		// we still want to get connection even if the cert is expired, or if
		// the hostname doesn't match
		InsecureSkipVerify: true,
	}

	// custom address parsing to allow default port and address override
	if !strings.Contains(host, ":") {
		host = host + ":443"
	} else if match := addrOverride.FindStringSubmatch(host); match != nil {
		config.ServerName = match[1]
		host = match[2]
	}

	// make a connection to get the certificate
	conn, err := tls.DialWithDialer(c.dialer, "tcp", host, &config)
	if err != nil {
		return
	}
	conn.Close()

	// check all certificates in the chain for expiration
	for _, cert := range conn.ConnectionState().PeerCertificates {
		if c.warn_if_expired_at.After(cert.NotAfter) {
			is_expired = true
			fmt.Printf("Certificate for %s (%s) expires in %s\n",
				host, cert.Subject.CommonName,
				humanize.Time(cert.NotAfter))
		}
	}

	// TODO: validate hostname and chain of trust

	return
}
