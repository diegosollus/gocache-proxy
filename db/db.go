package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB inicializa a conexão com o banco de dados SQLite e cria a tabela 'blocked_ips' se ela não existir.
// Ela configura algumas propriedades da conexão.
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "gocache-proxy.db")
	if err != nil {
		panic("Could not connect to database: " + err.Error())
	}

	DB.SetMaxOpenConns(10)
	DB.SetConnMaxIdleTime(5)

	if createTableBlockIps() {
		log.Printf("Table block_ips created successfully")
		insertBlockedIPs()
	}
}

// Carrega IPs bloqueados do banco de dados.
func LoadBlockedIPs(DB *sql.DB) ([]string, error) {
	var blockedIPs []string

	rows, err := DB.Query("SELECT ip_address FROM blocked_ips")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			return nil, err
		}
		blockedIPs = append(blockedIPs, ip)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return blockedIPs, nil
}

// createTableBlockIps cria a tabela 'blocked_ips' se ela não existir.
// Retorna true se a tabela acabou de ser criada.
func createTableBlockIps() bool {
	createIPTable := `
	CREATE TABLE IF NOT EXISTS blocked_ips (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip_address TEXT NOT NULL UNIQUE,
		blocked_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Verifica se a tabela já existe antes de criá-la.
	_, err := DB.Exec(createIPTable)
	if err != nil {
		panic("Could not create blocked_ips table: " + err.Error())
	}

	// Verifica se a tabela está vazia (ou seja, acabou de ser criada).
	var count int
	err = DB.QueryRow(`SELECT COUNT(*) FROM blocked_ips`).Scan(&count)
	if err != nil {
		log.Printf("Could not check if blocked_ips table is empty: %v", err)
		return false
	}

	return count == 0
}

// insertBlockedIPs insere 5 registros de IPs bloqueados na tabela 'blocked_ips' em uma única operação.
func insertBlockedIPs() {
	// Lista de IPs a serem bloqueados.
	blockedIPs := []string{
		"172.16.0.1",
		"172.16.0.2",
		"172.16.0.3",
		"172.16.0.4",
		"172.16.0.5",
	}

	insertQuery := `INSERT INTO blocked_ips (ip_address) VALUES `
	values := make([]interface{}, 0, len(blockedIPs))

	// Adiciona os placeholders de valores na query e armazena os IPs na lista de valores.
	for i, ip := range blockedIPs {
		if i > 0 {
			insertQuery += ", "
		}
		insertQuery += "(?)"
		values = append(values, ip)
	}

	_, err := DB.Exec(insertQuery, values...)
	if err != nil {
		log.Printf("Could not insert IPs: %v", err)
	} else {
		log.Printf("Inserted all IPs successfully")
	}
}
