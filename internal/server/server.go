package server

import (
	"fmt"
	"gocache-proxy/internal/configs"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Função principal que inicializa e executa o servidor.
func Run() error {
	// Carrega a configuração.
	config, err := configs.LoadConfig("data")
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)

	// Configura os recursos de proxy.
	for _, resource := range config.Resources {
		// Analisa e valida a URL de destino.
		url, err := url.Parse(resource.DestinationURL)
		if err != nil {
			log.Printf("Invalid URL for resource %s: %v", resource.Endpoint, err)
			continue // Pula este recurso, mas continua com os outros.
		}

		// Cria um proxy reverso para o recurso.
		proxy := NewProxy(url)
		handler := ProxyRequestHandler(proxy, url, resource.Endpoint)

		// Registra a rota.
		mux.HandleFunc(resource.Endpoint, handler)
	}

	// Configura o servidor HTTP com timeout.
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.Server.Host, config.Server.ListenPort),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Inicia o servidor.
	log.Printf("Starting server on %s:%s", config.Server.Host, config.Server.ListenPort)
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("could not start the server: %v", err)
	}

	return nil
}
