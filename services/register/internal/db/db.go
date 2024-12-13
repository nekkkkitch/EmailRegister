package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Config struct {
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	User     string `yaml:"user" env:"USER" env-default:"postgres"`
	Password string `yaml:"password" env:"password" env-default:"postgres"`
	DBName   string `yaml:"dbname" env:"DBNAME" env-default:"chat"`
}

type DB struct {
	config *Config
	db     *pgx.Conn
}

// Создает соединение с существующей БД
func New(cfg *Config) (*DB, error) {
	d := &DB{config: cfg}
	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err := pgx.Connect(context.Background(), connection)
	log.Println("Connecting to: " + connection)
	if err != nil {
		return nil, err
	}
	d.db = db
	return d, nil
}

func (d *DB) AddUser(email, password string) error {
	log.Println("Trying to insert user " + email)
	exists, err := d.CheckUserInDB(email)
	if err != nil {
		return err
	}
	if exists {
		log.Println("User already in db, still sending code")
		return nil
	}
	_, err = d.db.Exec(context.Background(), `insert into public.users(email, password) values($1, $2)`, email, password)
	if err != nil {
		log.Println("Cant insert user:", err)
		return err
	}
	log.Printf("User %v\n added successfully", email)
	return nil
}

func (d *DB) SetUserVerificationStatus(email string, status bool) error {
	_, err := d.db.Exec(context.Background(), `update public.users set verificated = $1 where email = $2`, status, email)
	if err != nil {
		log.Println("Cant set user verification status:", err)
		return err
	}
	return nil
}

func (d *DB) CheckUserInDB(email string) (bool, error) {
	var id pgtype.Int4
	err := d.db.QueryRow(context.Background(), `select id from public.users where email=$1`, email).Scan(&id)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return false, nil
		}
		log.Println("Fail checking user:", err)
		return false, err
	}
	return true, nil
}
