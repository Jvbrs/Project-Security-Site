package checks

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

func ExamineSSLSecurity(siteURL string) (bool, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}

	client := &http.Client{Transport: transport}

	resp, err := client.Head(siteURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.TLS != nil {
		currentTime := time.Now()
		for _, cert := range resp.TLS.PeerCertificates {
			if currentTime.Before(cert.NotBefore) {
				return false, fmt.Errorf("o certificado ainda não é válido")
			} else if currentTime.After(cert.NotAfter) {
				return false, fmt.Errorf("o certificado expirou")
			}
		}
		return true, nil
	}

	return false, fmt.Errorf("a conexão não utiliza SSL/TLS")
}

func ExaminePrivacyPolicy(siteURL string) (bool, error) {
	policyURLPatterns := []string{
		"/privacy",
		"/legal/privacy-policy",
		"/privacy-policy",
		"/terms",
		"/terms-of-service",
	}

	client := &http.Client{Timeout: 10 * time.Second}
	for _, pattern := range policyURLPatterns {
		policyURL := siteURL + pattern
		resp, err := client.Head(policyURL)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return true, nil
		}
	}

	return false, fmt.Errorf("a página de política de privacidade não está acessível")
}
