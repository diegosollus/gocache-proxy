package server

import (
	"fmt"
	utils "gocache-proxy/internal/utils"
	"gocache-proxy/security"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// NewProxy cria um novo proxy reverso para o host de destino.
func NewProxy(target *url.URL) *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(target)
}

// ProxyRequestHandler retorna um handler HTTP que processa a requisição e a redireciona ao backend.
func ProxyRequestHandler(proxy *httputil.ReverseProxy, targetURL *url.URL, endpoint string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now().UTC()

		// Exercício 01: Caso o IP esteja dentro de uma lista ele deve ser bloqueado, retornando status code 403.
		checkBlockedIP(w, r)

		// Exercício 02: Caso o subdomínio seja www, sobrescreva a URI para ser /site/www/<restante da uri>.
		handleSubdomain(r)

		// Exercício 03: Caso a query string ou a payload caia em um critério da regex, bloqueie a requisição retornando o status 403.
		checkMaliciousContent(w, r)

		fmt.Printf("[ GoCache Proxy ] Request received: %s -> Redirecting to: %s at %s\n", r.URL.Path, targetURL, startTime)

		r.URL.Host = targetURL.Host
		r.URL.Scheme = targetURL.Scheme
		r.Header.Set("X-Forwarded-Host", r.Host)
		r.Host = targetURL.Host
		r.URL.Path = strings.TrimPrefix(r.URL.Path, endpoint)

		fmt.Printf("[ GoCache Proxy ] Forwarding to: %s at %s\n", r.URL, startTime)

		// Redireciona a requisição para o servidor de backend via proxy.
		proxy.ServeHTTP(w, r)
	}
}

// Verifica se o IP está bloqueado.
func checkBlockedIP(w http.ResponseWriter, r *http.Request) {
	clientIP := utils.GetIPAddress(r)
	fmt.Printf("[ GoCache Proxy ] Request from IP: %s\n", clientIP)

	if utils.IsBlocked(clientIP) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		fmt.Printf("[ GoCache Proxy ] Blocked request from IP: %s\n", clientIP)
		return
	}
}

// Manipula o subdomínio "www".
func handleSubdomain(r *http.Request) {
	subdomain := utils.GetSubdomain(r)
	fmt.Printf("[ GoCache Proxy ] Subdomain detected: %s\n", subdomain)

	if subdomain == "www" {
		newPath := fmt.Sprintf("/site/www%s", r.URL.Path)
		fmt.Printf("[ GoCache Proxy ] Rewriting URI to: %s\n", newPath)
		r.URL.Path = newPath
	}
}

// Verifica se a query string ou payload contém tags HTML maliciosas.
func checkMaliciousContent(w http.ResponseWriter, r *http.Request) {
	// Verificação da query string.
	escapedQuery, err := url.QueryUnescape(r.URL.RawQuery)
	if err != nil {
		fmt.Printf("[ GoCache Proxy ] Error unescaping query: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Printf("[ GoCache Proxy ] Unescaped query: %s\n", escapedQuery)

	if security.TagRegex.MatchString(escapedQuery) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		fmt.Printf("[ GoCache Proxy ] Blocked request with malicious query string\n")
		return
	}

	// Verificação do body (para requisições POST ou PUT)
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("[ GoCache Proxy ] Error reading request body: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		bodyString := string(bodyBytes)

		// Reinicializa o body para permitir que ele seja lido pelo proxy depois.
		r.Body = io.NopCloser(strings.NewReader(bodyString))

		if security.TagRegex.MatchString(bodyString) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			fmt.Printf("[ GoCache Proxy ] Blocked request with malicious payload\n")
			return
		}
	}
}
