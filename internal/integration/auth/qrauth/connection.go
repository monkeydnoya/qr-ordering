package qrauth

import (
	"context"
	"net/http"

	"qr-ordering-service/internal/integration/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type QrAuthConfig struct {
	// Host string
	// Port string
	Url string
}

// TODO: NewConnection to QRAuth service that validate tables.
type qrAuthConnection struct {
	conf   QrAuthConfig
	logger logx.Logger
	client http.Client
}

func (qrAuthConf *QrAuthConfig) NewConnection() (auth.QrAuth, error) {
	qrClient := &http.Client{}
	auth := qrAuthConnection{
		conf:   *qrAuthConf,
		logger: logx.WithContext(context.Background()),
		client: *qrClient,
	}
	return &auth, nil
}
