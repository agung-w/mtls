package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Inisiasi Echo
	e := echo.New()

	// Tambahkan middleware umum
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Load CA certificate
	caCert, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		e.Logger.Fatalf("Error reading CA certificate: %v", err)
	}

	// Create certificate pool dan tambahkan CA certificate
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		e.Logger.Fatal("Failed to append CA certificate to pool")
	}

	// Load server certificate dan key
	serverCert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		e.Logger.Fatalf("Error loading server certificate: %v", err)
	}

	// Rute-rute API
	e.GET("/", func(c echo.Context) error {
		// Ambil informasi sertifikat klien
		req := c.Request()

		// Periksa apakah request menggunakan TLS
		if req.TLS == nil {
			return c.String(http.StatusInternalServerError, "Tidak menggunakan koneksi TLS")
		}

		// Dapatkan sertifikat klien
		peerCerts := req.TLS.PeerCertificates
		var clientCN string
		if len(peerCerts) > 0 {
			clientCN = peerCerts[0].Subject.CommonName
			e.Logger.Printf("Authenticated client: %s", clientCN)
		}

		// Contoh response
		return c.String(http.StatusOK, fmt.Sprintf("Selamat datang, %s! Koneksi mTLS berhasil.", clientCN))
	})

	// Endpoint contoh lainnya
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	// Konfigurasi server
	s := &http.Server{
		Addr: ":8443", // Port HTTPS
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{serverCert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    caCertPool,
			MinVersion:   tls.VersionTLS12,
		},
	}

	// Start server
	e.Logger.Fatal(e.StartServer(s))
}
