package postgres

import (
	"context"
	"qr-ordering-service/internal/data/database"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
}

type pgConnection struct {
	conf PostgresConfig
	log  logx.Logger
	db   *gorm.DB
}

func (pConf *PostgresConfig) NewConnection() (database.Database, error) {
	pg := pgConnection{
		conf: *pConf,
		log:  logx.WithContext(context.Background()),
	}

	return &pg, nil
}
