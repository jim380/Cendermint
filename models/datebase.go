package models

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jim380/Cendermint/migrations"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

/*
Open a db connection with Postgres;
The caller is responsible for closing the connection
*/
func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	return db, nil
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "admin",
		Password: "1234567890",
		Database: "cendermint",
		SSLMode:  "disable",
	}
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	// In case the dir is an empty string, they probably meant the current directory and goose wants a period for that.
	if dir == "" {
		dir = "."
	}
	goose.SetBaseFS(migrationsFS)
	defer func() {
		// Ensure that we remove the FS on the off chance some other part of our app uses goose for migrations and doesn't want to use our FS.
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func SetupDatabase() *sql.DB {
	dbConfig := DefaultPostgresConfig()
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Using db config", dbConfig.String()))
	db, err := Open(dbConfig)
	if err != nil {
		zap.L().Fatal("\t", zap.Bool("Success", false), zap.String("Database connection", "failed with error:"+err.Error()))
	} else {
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Database connection", "ok"))
	}

	return db
}

func MigrateDatabase(db *sql.DB) {
	// migrate general tables
	err := MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	// migrate chain specific tables
}
