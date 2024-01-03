package main

import (
	"bytes"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

var googleTimeout = 5 * time.Second

func TestCheckHostOneOneOneOne(t *testing.T) {
	checker := NewSimpleHostChecker(
		&net.Dialer{Timeout: googleTimeout},
		time.Now().AddDate(0, 0, 7),
	)
	expired, err := checker.CheckHost("one.one.one.one:1.1.1.1:443")
	if err != nil {
		t.Errorf("checkHost() with valid host returned an error: %v", err)
	}
	if expired {
		t.Errorf("checkHost() with valid host reported expired certificate")
	}
}

func TestCheckHostGoogleValid(t *testing.T) {
	checker := NewSimpleHostChecker(
		&net.Dialer{Timeout: googleTimeout},
		time.Now().AddDate(0, 0, 7),
	)
	expired, err := checker.CheckHost("google.com")
	if err != nil {
		t.Errorf("checkHost() with valid host returned an error: %v", err)
	}
	if expired {
		t.Errorf("checkHost() with valid host reported expired certificate")
	}
}

func captureOutput(f func()) string {
	// Keep backup of the real stdout
	old := os.Stdout
	// Create a pipe for capturing output
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the function
	f()

	// Close the pipe
	w.Close()
	// Restore the original stdout
	os.Stdout = old

	var buf bytes.Buffer
	// Read the captured output
	buf.ReadFrom(r)
	return buf.String()
}

func TestGoogleCertInFiveYears(t *testing.T) {
	checker := NewSimpleHostChecker(
		&net.Dialer{Timeout: googleTimeout},
		// Set the warning threshold to 5 years from now
		time.Now().AddDate(5, 0, 0),
	)
	var isExpired bool
	var err error
	output := captureOutput(func() {
		isExpired, err = checker.CheckHost("google.com")
	})
	if err != nil {
		t.Fatalf("checkHost() returned an error: %v", err)
	}
	if !strings.Contains(output, "Certificate for google.com:443 (*.google.com) expires in") {
		t.Errorf("Expected output to contain certificate expiration warning, got %q", output)
	}
	if !isExpired {
		t.Errorf("Google's certificate is not marked as expiring within 5 years")
	}
}

func TestCheckHostUnreachable(t *testing.T) {
	checker := NewSimpleHostChecker(
		&net.Dialer{Timeout: googleTimeout},
		time.Now().AddDate(0, 0, 7),
	)
	isExpired, err := checker.CheckHost("some-unreachable-domain")
	if err == nil {
		t.Errorf("checkHost() with unreachable host did not return an error")
	}
	if isExpired {
		t.Errorf("checkHost() with unreachable host reported expired certificate")
	}
}

type MockHostChecker struct {
	MockCheckHost func(host string) (bool, error)
}

func (m *MockHostChecker) CheckHost(host string) (bool, error) {
	return m.MockCheckHost(host)
}

func TestRunChecks(t *testing.T) {

	// Mock checkHost function
	mockCheckHost := func(host string) (bool, error) {
		if host == "test.com" {
			return false, nil // test.com has a valid certificate
		}
		// unreachable domains error out
		if strings.HasPrefix(host, "some-unreachable-domain") {
			return false, &net.OpError{
				Op:  "dial",
				Net: "tcp",
				Err: &net.DNSError{Err: "no such host"},
			}
		}
		return true, nil // other domains are considered expired
	}

	checker := &MockHostChecker{MockCheckHost: mockCheckHost}

	// Define test cases
	testCases := []struct {
		name           string
		endpoints      []string
		expectedCode   int
		expectedOutput string
	}{
		{"All Valid", []string{"test.com", "test.com"}, 0, ""},
		{"One Expired", []string{"test.com", "expired.com"}, 1, ""},
		{"Unreachable", []string{"some-unreachable-domain"}, 1, "can't check"},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var exitCode int
			output := captureOutput(func() {
				exitCode = RunChecks(checker, tc.endpoints, 2)
			})
			if exitCode != tc.expectedCode {
				t.Errorf("Expected exit code %d, got %d", tc.expectedCode, exitCode)
			}
			if tc.expectedOutput != "" && !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("Expected output %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}
