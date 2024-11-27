package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// Load CA Certificate
	caCert, err := ioutil.ReadFile("/etc/ssl/certs/vg-ca.pem")
	if err != nil {
		log.Fatalf("Gagal membaca CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Load Client Certificate dan Key
	cert, err := tls.LoadX509KeyPair("/etc/ssl/certs/client_cert.pem", "/etc/ssl/private/client_key.key")
	if err != nil {
		log.Fatalf("Gagal membaca client certificate: %v", err)
	}

	// Konfigurasi TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	// Buat transport khusus dengan konfigurasi TLS
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	// Buat client HTTP dengan transport khusus
	client := &http.Client{
		Transport: transport,
	}

	// Lakukan request ke server
	resp, err := client.Get("https://agung-w-system-product-name:8443")
	if err != nil {
		log.Fatalf("Gagal melakukan request: %v", err)
	}
	defer resp.Body.Close()

	// Baca response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Gagal membaca response: %v", err)
	}

	fmt.Println(string(body))
}
