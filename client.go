package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func createMTLSClient() *http.Client {
	// Load CA certificate
	caCert, err := os.ReadFile("/etc/ssl/certs/vg-ca.pem")
	if err != nil {
		log.Fatalf("Error reading CA certificate: %v", err)
	}

	// Create certificate pool and add CA certificate
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("Failed to append CA certificate to pool")
	}

	// Load client certificate and key
	clientCert, err := tls.LoadX509KeyPair("/etc/ssl/certs/client_cert.pem", "/etc/ssl/private/client_key.key")
	if err != nil {
		log.Fatalf("Error loading client certificate: %v", err)
	}

	// Configure TLS with mutual authentication
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12, // Enforce minimum TLS version
	}

	// Create HTTP client with mTLS configuration and timeouts
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		},
		Timeout: 10 * time.Second, // Overall request timeout
	}
}

func makeAPIRequest(client *http.Client, url string) {
	// Perform GET request
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Print response details
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}

func main() {
	// Create mTLS HTTP client
	mtlsClient := createMTLSClient()

	// API endpoint (replace with your actual endpoint)
	apiURL := "https://your-mtls-api-endpoint.com/api/path"

	// Make API request
	makeAPIRequest(mtlsClient, apiURL)
}
