package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Model interface {
	StoreBirthday(string, time.Time) error
	RetrieveBirthdayMessage(string) (string, error)
}

type DB struct {
	pool *pgxpool.Pool
	config *pgxpool.Config
	ctx context.Context
}

func GetConnection(url string) (*DB, error) {
	var err error
	db := DB{}
	db.ctx = context.Background()

	db.config, err = pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	db.config.MaxConns = 20
	db.config.MaxConnLifetime = 10 * time.Second

	db.pool, err = pgxpool.ConnectConfig(db.ctx, db.config)
	if err != nil {
		return nil, err
	}
	return &db, nil
}

func (db *DB) CloseConnections() {
	db.pool.Close()
}

func (db *DB) StoreBirthday(name string, dateOfBirth time.Time) error {
	_, err := db.pool.Exec(db.ctx, "SELECT hello.store_birthday($1::text, $2::date)", name, dateOfBirth)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) RetrieveBirthdayMessage(name string) (message string, err error) {
	var result *string
	row := db.pool.QueryRow(db.ctx, "SELECT hello.retrieve_birthday_message($1::text, $2::date)", name, time.Now())
	err = row.Scan(&result)
	if err != nil {
		return
	}
	message = *result
	return
}

