package main

import (
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/caddyserver/certmagic"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func initCertMagic() {
	certmagic.DefaultACME.Agreed = true
	email := getEnv("CERTMAGIC_EMAIL", "")
	if email != "" {
		certmagic.DefaultACME.Email = email
	}
	isProd := getEnv("ENV", "false")
	if isProd == "true" {
		certmagic.DefaultACME.CA = certmagic.LetsEncryptProductionCA
	} else {
		certmagic.DefaultACME.CA = certmagic.LetsEncryptStagingCA
	}
}

func main() {
	initCertMagic()
	// Start a reverse proxy to the real site
	target := "api.openai.com"
	director := func(req *http.Request) {
		req.URL.Scheme = "https"
		req.URL.Host = target
		req.Host = target
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
	mux := http.NewServeMux()
	mux.Handle("/v1/", proxyHandler)
	myDomain := "localhost"
	certmagic.HTTPS([]string{myDomain}, mux)
}
