package svc

import (
	"context"
	"qr-ordering-service/internal/config"
	"qr-ordering-service/internal/data/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	log    logx.Logger
	Config config.Config
	Db     database.Database
}

func NewServiceContext(c config.Config) *ServiceContext {
	var database database.Database
	var err error

	SC := ServiceContext{
		Config: c,
		log:    logx.WithContext(context.Background()).WithFields(logx.Field("service", "ordering")),
	}
	switch c.Database.DatabaseType {
	case "postgres":
		database, err = c.Database.Postgres.NewConnection()
		if err != nil {
			SC.log.Errorf("postgres database: %s", err)
			panic(err)
		}
		SC.log.Infof("postgres database: connected successfully: %s:%s", c.Database.Postgres.Host, c.Database.Postgres.Port)
	default:
		panic("unknown database type: " + c.Database.DatabaseType)
	}
	if err := database.Start(); err != nil {
		SC.log.Errorf("postgres database: could not start connection: %s", err)
		panic(err)
	}
	SC.Db = database
	return &SC
}
