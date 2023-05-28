package config

import (
	"qr-ordering-service/internal/data/database/postgres"
	"qr-ordering-service/internal/integration/auth/qrauth"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Database
	QrAuth
}

type Database struct {
	DatabaseType string
	Postgres     postgres.PostgresConfig
}

type QrAuth struct {
	QrAuthType string
	Qrauth     qrauth.QrAuthConfig
}
