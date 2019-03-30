package services

const (
	Postgres = "postgres"
)

type Service interface {
	Name() string
	Start() error
	Stop() error
	Started() bool
}
