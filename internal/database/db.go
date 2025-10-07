package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "1234")
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "smartpicks")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro ao conectar no MySQL:", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)

	if err = DB.Ping(); err != nil {
		log.Fatal("MySQL não respondeu:", err)
	}
	log.Println("Conectado ao MySQL com sucesso!")

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
	var tableName string
	if err := db.QueryRow("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users'").Scan(&tableName); err != nil {
		return nil
	}

	var dataType string
	err := db.QueryRow(`
		SELECT DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users' AND COLUMN_NAME = 'avatar'
	`).Scan(&dataType)

	if err == sql.ErrNoRows {
		if _, alterErr := db.Exec("ALTER TABLE users ADD COLUMN avatar MEDIUMTEXT NULL DEFAULT NULL"); alterErr != nil {
			return fmt.Errorf("falha ao adicionar coluna avatar: %w", alterErr)
		}
		log.Println("[DB] Coluna 'avatar' criada como MEDIUMTEXT")
		return nil
	} else if err != nil {
		return err
	}
	typ := strings.ToLower(dataType)
	if typ != "mediumtext" && typ != "longtext" {
		idxRows, idxErr := db.Query(`
			SELECT DISTINCT INDEX_NAME
			FROM INFORMATION_SCHEMA.STATISTICS
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users' AND COLUMN_NAME = 'avatar'
		`)
		if idxErr == nil {
			defer idxRows.Close()
			for idxRows.Next() {
				var idxName string
				if scanErr := idxRows.Scan(&idxName); scanErr == nil {
					if strings.ToUpper(idxName) == "PRIMARY" {
						continue
					}
					if _, dropErr := db.Exec("ALTER TABLE users DROP INDEX " + idxName); dropErr != nil {
						log.Printf("[DB] Aviso: falha ao remover índice %s: %v", idxName, dropErr)
					} else {
						log.Printf("[DB] Índice removido: %s", idxName)
					}
				}
			}
		}

		if _, alterErr := db.Exec("ALTER TABLE users MODIFY COLUMN avatar MEDIUMTEXT NULL"); alterErr != nil {
			return fmt.Errorf("falha ao alterar tipo da coluna avatar para MEDIUMTEXT: %w", alterErr)
		}
		log.Println("[DB] Coluna 'avatar' atualizada para MEDIUMTEXT")
	}
	return nil
}
