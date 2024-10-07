package main

import (
	"gocache-proxy/db"
	"gocache-proxy/internal/server"
	"gocache-proxy/internal/utils"
	"log"
)

func main() {
	// Inicializa o banco de dados.
	db.InitDB()

	// Carrega os IPs bloqueados.
	utils.LoadBlockedIPs()

	// Inicia o servidor.
	if err := server.Run(); err != nil {
		log.Fatalf("could not start the server: %v", err)
	}
}
