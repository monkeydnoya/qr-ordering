package svc

import (
	"context"
	"qr-ordering-service/internal/config"
	"qr-ordering-service/internal/data/database"
	"qr-ordering-service/internal/integration/auth"
	"qr-ordering-service/internal/middleware"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	log        logx.Logger
	Config     config.Config
	Middleware rest.Middleware

	Db     database.Database
	QrAuth auth.QrAuth
}

func NewServiceContext(c config.Config) *ServiceContext {
	var database database.Database
	var qrauth auth.QrAuth
	var err error

	SC := ServiceContext{
		Config: c,
		log:    logx.WithContext(context.Background()).WithFields(logx.Field("service", "ordering")),
	}
	switch c.QrAuth.QrAuthType {
	case "qrauth":
		qrauth, err = c.QrAuth.Qrauth.NewConnection()
		if err != nil {
			SC.log.Errorf("qr auth: %s", err)
			panic(err)
		}
		SC.log.Infof("qr auth: connected successfully")
	default:
		panic("unknown qrauth type: " + c.QrAuth.QrAuthType)
	}
	SC.QrAuth = qrauth
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
	SC.Middleware = middleware.NewMiddleware(SC.QrAuth).Handler
	return &SC
}
