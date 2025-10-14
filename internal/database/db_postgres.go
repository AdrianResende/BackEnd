package database

import (
"database/sql"
"fmt"
"log"
"os"

_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
// Para Neon PostgreSQL, use a DATABASE_URL diretamente
databaseURL := os.Getenv("DATABASE_URL")
if databaseURL == "" {
// Fallback para variáveis separadas
dbUser := getEnv("DB_USER", "postgres")
dbPassword := getEnv("DB_PASSWORD", "")
dbHost := getEnv("DB_HOST", "localhost")
dbPort := getEnv("DB_PORT", "5432")
dbName := getEnv("DB_NAME", "smartpicks")
databaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", 
dbUser, dbPassword, dbHost, dbPort, dbName)
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
log.Println("Conectado ao PostgreSQL com sucesso!")

if err := ensureSchema(DB); err != nil {
log.Printf("[DB] Aviso: falha ao garantir schema: %v", err)
}
}

func getEnv(key, defaultValue string) string {
if value := os.Getenv(key); value != "" {
return value
}
return defaultValue
}

func ensureSchema(db *sql.DB) error {
// Verificar se a tabela users existe
var exists bool
err := db.QueryRow(`
SELECT EXISTS (
SELECT FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name = 'users'
)
`).Scan(&exists)

if err != nil {
return err
}

if !exists {
log.Println("[DB] Tabela users não existe, será criada no primeiro uso")
return nil
}

// Verificar se a coluna avatar existe
var columnExists bool
err = db.QueryRow(`
SELECT EXISTS (
SELECT FROM information_schema.columns 
WHERE table_name = 'users' 
AND column_name = 'avatar'
)
`).Scan(&columnExists)

if err != nil {
return err
}

if !columnExists {
if _, alterErr := db.Exec("ALTER TABLE users ADD COLUMN avatar TEXT NULL DEFAULT NULL"); alterErr != nil {
return fmt.Errorf("falha ao adicionar coluna avatar: %w", alterErr)
}
log.Println("[DB] Coluna 'avatar' criada como TEXT")
}

return nil
}
