package db

import (
	"auth-service/pkg/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBUser, cfg.DBPass, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (db *Database) SaveRefreshToken(userID, token string) error {
	_, err := db.Exec("INSERT INTO refresh_tokens (user_id, token) VALUES ($1, $2)", userID, token)
	return err
}

func (db *Database) GetUserIDByRefreshToken(token string, userID *string) error {
	return db.QueryRow("SELECT user_id FROM refresh_tokens WHERE token = $1", token).Scan(userID)
}

func (db *Database) DeleteRefreshToken(userID, token string) error {
	_, err := db.Exec("DELETE FROM refresh_tokens WHERE user_id = $1 AND token = $2", userID, token)
	return err
}

func (db *Database) GetUserIPAddress(userID string) (string, error) {
	var ipAddress string
	err := db.QueryRow("SELECT ip_address FROM users WHERE id = $1", userID).Scan(&ipAddress)
	return ipAddress, err
}
