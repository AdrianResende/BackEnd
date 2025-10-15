package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	// Detectar ambiente baseado na presença de DATABASE_URL
	// DATABASE_URL presente: PostgreSQL (pode ser local ou Vercel/Neon)
	// DATABASE_URL ausente: MySQL (desenvolvimento local)
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL != "" {
		// Usar PostgreSQL
		connectPostgreSQL(databaseURL)
	} else {
		// Usar MySQL
		connectMySQL()
	}
}

func connectPostgreSQL(databaseURL string) {
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

	// Detectar se é local ou produção pela URL
	ambiente := "Local"
	if !strings.Contains(databaseURL, "localhost") && !strings.Contains(databaseURL, "127.0.0.1") {
		ambiente = "Produção (Vercel/Neon)"
	}

	log.Printf("✓ Conectado ao PostgreSQL (%s)", ambiente)
}

func connectMySQL() {
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "1234")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "smartpicks")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro ao conectar no MySQL:", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	if err = DB.Ping(); err != nil {
		log.Fatal("MySQL não respondeu:", err)
	}
	log.Println("✓ Conectado ao MySQL (Desenvolvimento Local)")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
