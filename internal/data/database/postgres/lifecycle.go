package postgres

import (
	"fmt"
	"net/url"
	"qr-ordering-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (pg *pgConnection) Start() error {
	var err error
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		url.QueryEscape(pg.conf.Username),
		url.QueryEscape(pg.conf.Password),
		url.QueryEscape(pg.conf.Host),
		pg.conf.Port,
		url.QueryEscape(pg.conf.DbName),
	)
	pg.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		pg.log.Errorw("postgres: failed to create new posgresql connection",
			logx.LogField{Key: "err", Value: err.Error()})
		return err
	}

	if err = pg.db.AutoMigrate(&types.OrderEntity{}); err != nil {
		pg.log.Errorw("postgres: failed to migrate order entity",
			logx.LogField{Key: "err", Value: err.Error()})
		return err
	}

	if err = pg.db.AutoMigrate(&types.ItemEntity{}); err != nil {
		pg.log.Errorw("postgres: failed to migrate item entity",
			logx.LogField{Key: "err", Value: err.Error()})
		return err
	}
	return nil
}

func (pg *pgConnection) Stop() error {
	dbInstance, err := pg.db.DB()
	if err == nil {
		dbInstance.Close()
	}
	return nil
}
