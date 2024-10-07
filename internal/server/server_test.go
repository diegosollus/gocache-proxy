package server_test

import (
	"fmt"
	"gocache-proxy/internal/server"
	"gocache-proxy/internal/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// Carrega os IPs bloqueados no teste.
func loadBlockedIPs() {
	utils.BlockedIPs = []string{
		"172.16.0.1",
		"172.16.0.2",
		"172.16.0.3",
		"172.16.0.4",
		"172.16.0.5",
	}
}

// TestProxyBlockedIP testa se o handler bloqueia corretamente IPs da lista de bloqueados.
func TestProxyBlockedIP(t *testing.T) {
	loadBlockedIPs()

	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Request received by backend")
	}))
	defer targetServer.Close()

	targetURL, _ := url.Parse(targetServer.URL)

	// Cria o proxy
	proxy := server.NewProxy(targetURL)

	// Testa requisição de IP bloqueado
	req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
	req.RemoteAddr = "172.16.0.1:12345" // Simula um IP bloqueado

	w := httptest.NewRecorder()

	handler := server.ProxyRequestHandler(proxy, targetURL, "/test")
	handler(w, req)

	if w.Result().StatusCode != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Result().StatusCode)
	}

	body, _ := io.ReadAll(w.Result().Body)
	if !strings.Contains(string(body), "Forbidden") {
		t.Errorf("Expected Forbidden response body, got %s", string(body))
	}
}

// TestProxyAllowedIP testa se o handler permite requisições de IPs não bloqueados.
func TestProxyAllowedIP(t *testing.T) {
	loadBlockedIPs()

	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Request received by backend")
	}))
	defer targetServer.Close()

	targetURL, _ := url.Parse(targetServer.URL)

	// Cria o proxy.
	proxy := server.NewProxy(targetURL)

	// Testa requisição de IP permitido.
	req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
	req.RemoteAddr = "192.168.1.1:12345" // Simula um IP não bloqueado

	w := httptest.NewRecorder()

	handler := server.ProxyRequestHandler(proxy, targetURL, "/test")
	handler(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Result().StatusCode)
	}

	body, _ := io.ReadAll(w.Result().Body)
	if !strings.Contains(string(body), "Request received by backend") {
		t.Errorf("Expected backend response body, got %s", string(body))
	}
}

// TestProxyRewriteSubdomainWorldWideWeb testa se o handler sobrescreve a URI corretamente quando o subdomínio é "www".
func TestProxyRewriteSubdomainWorldWideWeb(t *testing.T) {
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Request received by backend")
	}))
	defer targetServer.Close()

	targetURL, _ := url.Parse(targetServer.URL)

	// Cria o proxy.
	proxy := server.NewProxy(targetURL)

	// Simula uma requisição com subdomínio "www".
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/test", nil)
	req.RemoteAddr = "192.168.1.1:12345" // Simula um IP não bloqueado.

	w := httptest.NewRecorder()

	handler := server.ProxyRequestHandler(proxy, targetURL, "/test")
	handler(w, req)

	// Verifica se a URI foi sobrescrita corretamente.
	expectedPath := "/site/www/test"
	if req.URL.Path != expectedPath {
		t.Errorf("Expected rewritten path %s, got %s", expectedPath, req.URL.Path)
	}

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Result().StatusCode)
	}

	body, _ := io.ReadAll(w.Result().Body)
	if !strings.Contains(string(body), "Request received by backend") {
		t.Errorf("Expected backend response body, got %s", string(body))
	}
}

// TestProxyMaliciousQueryString testa se o handler bloqueia requisições com query string maliciosa.
func TestProxyMaliciousQueryString(t *testing.T) {
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Request received by backend")
	}))
	defer targetServer.Close()

	targetURL, _ := url.Parse(targetServer.URL)

	// Cria o proxy
	proxy := server.NewProxy(targetURL)

	// Simula uma requisição com query string maliciosa.
	req := httptest.NewRequest(http.MethodGet, "http://example.com/test?<script>alert(1)</script>", nil)
	req.RemoteAddr = "192.168.1.1:12345" // Simula um IP não bloqueado.

	w := httptest.NewRecorder()

	handler := server.ProxyRequestHandler(proxy, targetURL, "/test")
	handler(w, req)

	if w.Result().StatusCode != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Result().StatusCode)
	}

	body, _ := io.ReadAll(w.Result().Body)
	if !strings.Contains(string(body), "Forbidden") {
		t.Errorf("Expected Forbidden response body, got %s", string(body))
	}
}
