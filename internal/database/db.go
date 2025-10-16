package database

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL não definida. Configure a string de conexão do PostgreSQL no .env")
	}

	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Erro ao conectar no PostgreSQL:", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)

	if err = DB.Ping(); err != nil {
		log.Fatal("PostgreSQL não respondeu:", err)
	}

	ambiente := "Local"
	if !strings.Contains(databaseURL, "localhost") && !strings.Contains(databaseURL, "127.0.0.1") {
		ambiente = "Produção (Vercel/Neon)"
	}

	log.Printf("✓ Conectado ao PostgreSQL (%s)", ambiente)
}
