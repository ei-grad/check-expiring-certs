package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
)

func main() {

	exitcode := 0

	// get timeout from opts
	timeout := flag.Duration("timeout", 2*time.Second, "timeout for connection")
	concurrency := flag.Int("concurrency", 128, "number of concurrent checks")
	flag.Parse()

	// endpoints to check
	endpoints := flag.Args()

	// semaphore to limit concurrency to a reasonable number
	semaphore := make(chan struct{}, *concurrency)

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

			is_expired, err := checkHost(i, *timeout)
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

func checkHost(host string, timeout time.Duration) (is_expired bool, err error) {
	dialer := &net.Dialer{
		Timeout: timeout,
	}
	conn, err := tls.DialWithDialer(dialer, "tcp", host, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return
	}
	conn.Close()
	for _, cert := range conn.ConnectionState().PeerCertificates {
		if time.Now().AddDate(0, 0, 7).After(cert.NotAfter) {
			is_expired = true
			fmt.Printf("Certificate for %s (%s) expires in %s\n",
				host, cert.Subject.CommonName,
				humanize.Time(cert.NotAfter))
		}
	}
	return
}
