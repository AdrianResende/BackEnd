package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	// TODO: Implementar carregamento de variáveis de ambiente
	// Por enquanto, usando configuração padrão
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "1234")
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "smartpicks")

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro ao conectar no MySQL:", err)
	}

	// Configurações de conexão
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)

	if err = DB.Ping(); err != nil {
		log.Fatal("MySQL não respondeu:", err)
	}
	log.Println("Conectado ao MySQL com sucesso!")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
