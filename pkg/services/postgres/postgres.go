package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq" // Postgres driver
)

type Service struct {
	config *Config
	db     *sql.DB
}

type Config struct {
	Addr            string        `envconfig:"addr" default:"postgres://expected:expected@localhost/expected?sslmode=disable"`
	ConnMaxLifetime time.Duration `envconfig:"connmaxlifetime" default:"5m"`
	MaxIdleConns    int           `envconfig:"maxidleconns" default:"4"`
	MaxOpenConns    int           `envconfig:"maxopenconns" default:"100"`
}

func New(config *Config) *Service {
	return &Service{
		config: config,
		db:     nil,
	}
}

func NewFromEnv() *Service {
	config := &Config{}
	if err := envconfig.Process("POSTGRES", config); err != nil {
		logrus.WithError(err).Fatalln("unable to process environment configuration")
	}
	return New(config)
}

func (srv *Service) Name() string {
	return "postgres"
}

func (srv *Service) Start() error {
	db, err := sql.Open("postgres", srv.config.Addr)
	if err != nil {
		return nil
	}
	db.SetConnMaxLifetime(srv.config.ConnMaxLifetime)
	db.SetMaxIdleConns(srv.config.MaxIdleConns)
	db.SetMaxOpenConns(srv.config.MaxOpenConns)
	srv.db = db
	return db.Ping()
}

func (srv *Service) Stop() error {
	return srv.db.Close()
}

func (srv *Service) Started() bool {
	return srv.db != nil && srv.db.Ping() == nil
}

func (srv *Service) Client() *sql.DB {
	return srv.db
}

func (srv *Service) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if srv.db == nil {
		return nil, errors.New("postgres database is not started")
	}
	return srv.db.ExecContext(ctx, query, args...)
}

func (srv *Service) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if srv.db == nil {
		return nil, errors.New("postgres database is not started")
	}
	return srv.db.QueryContext(ctx, query, args...)
}
