// server/server.go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Koneksi mTLS berhasil!\n")
}

func main() {
	// Load CA Certificate
	caCert, err := ioutil.ReadFile("/etc/ssl/certs/vg-ca.pem")
	if err != nil {
		log.Fatalf("Gagal membaca CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Konfigurasi TLS
	tlsConfig := &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caCertPool,
	}

	// Setup server dengan mTLS
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
		Handler:   http.HandlerFunc(helloHandler),
	}

	// Jalankan server dengan sertifikat server
	log.Println("Server mTLS berjalan di https://agung-w-system-product-name:8443")
	err = server.ListenAndServeTLS("/etc/ssl/certs/client_cert.pem", "/etc/ssl/private/client_key.key")
	if err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
