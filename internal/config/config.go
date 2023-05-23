package config

import (
	"qr-ordering-service/internal/data/database/postgres"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Database
}

type Database struct {
	DatabaseType string
	Postgres     postgres.PostgresConfig
}
