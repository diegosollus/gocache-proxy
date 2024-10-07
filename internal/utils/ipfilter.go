package utils

import (
	"gocache-proxy/db"
	"log"
	"net/http"
	"strings"
)

// Lista de IPs bloqueados (inicialmente vazia).
var BlockedIPs []string

// LoadBlockedIPs carrega os IPs bloqueados e preenche BlockedIPs.
func LoadBlockedIPs() {
	blockedIPs, err := db.LoadBlockedIPs(db.DB)

	if err != nil {
		log.Fatalf("Could not load blocked IPs: %v", err)
	}

	BlockedIPs = blockedIPs
}

// IsBlocked verifica se o IP está na lista de bloqueados.
func IsBlocked(ip string) bool {
	for _, blockedIP := range BlockedIPs {
		if ip == blockedIP {
			return true
		}
	}

	return false
}

// GetIPAddress recupera o IP do cliente.
func GetIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	// Remove a porta caso exista.
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}

	return ip
}
