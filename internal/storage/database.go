package storage

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	DB PgxConnInterface
}

func NewDB(connStr string) (*Database, error) {
	var conn *pgx.Conn
	var err error
	retries := 10

	for i := range retries {
		conn, err = pgx.Connect(context.Background(), connStr)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to DB %d/%d: %v", i+1, retries, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	db := &Database{DB: conn}
	_, err = db.DB.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS url_table (id SERIAL PRIMARY KEY,key TEXT NOT NULL UNIQUE,value TEXT NOT NULL UNIQUE)`)

	if err != nil {
		db.DB.Close(context.Background())
		return nil, err
	}
	return db, nil
}

func (d *Database) Set(key, value string) error {
	_, err := d.DB.Exec(
		context.Background(),
		`INSERT INTO url_table (key,value) values ($1,$2)`, key, value)
	return err
}

func (d *Database) Get(key string) (string, error) {
	var value string
	err := d.DB.QueryRow(context.Background(), `SELECT value FROM url_table WHERE key=$1`, key).Scan(&value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("Key dosn't exists in database")
		}
		return "", err
	}
	return value, nil
}
func (d *Database) Close() error {
	return d.DB.Close(context.Background())
}
