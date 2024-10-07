package utils

import (
	"net/http"
	"strings"
)

// GetSubdomain extrai o subdomínio da requisição.
func GetSubdomain(r *http.Request) string {
	parts := strings.Split(r.Host, ".")

	if len(parts) > 2 {
		return parts[0]
	}

	return ""
}
