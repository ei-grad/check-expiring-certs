package main

import "os"
import "time"
import "fmt"
import "crypto/tls"

import "github.com/dustin/go-humanize"

var exitcode int

func main() {
	for _, i := range os.Args[1:] {
		err := checkHost(i)
		if err != nil {
			fmt.Printf("can't check %s: %s\n", i, err)
			exitcode = 1
		}
	}
	os.Exit(exitcode)
}

func checkHost(host string) (err error) {
	conn, err := tls.Dial("tcp", host, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return
	}
	conn.Close()
	for _, cert := range conn.ConnectionState().PeerCertificates {
		if time.Now().AddDate(0, 0, 7).After(cert.NotAfter) {
			exitcode = 1
			fmt.Printf("Certificate for %s (%s) expires in %s\n",
				host, cert.Subject.CommonName,
				humanize.Time(cert.NotAfter))
		}
	}
	return
}
