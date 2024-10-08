package main

import (
	"gocache-proxy/db"
	"gocache-proxy/internal/httphelper"
	"gocache-proxy/internal/server"
	"log"
)

func main() {
	// Inicializa o banco de dados.
	db.InitDB()

	// Carrega os IPs bloqueados.
	httphelper.LoadBlockedIPs()

	// Inicia o servidor.
	if err := server.Run(); err != nil {
		log.Fatalf("could not start the server: %v", err)
	}
}
