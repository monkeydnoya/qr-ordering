package qrauth

import (
	"fmt"
	"net/http"
	"qr-ordering-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// Request exaple: http://localhost:8000/api/aut/table-validate/:table
func (qr *qrAuthConnection) ValidateTable(table types.Table) error {
	url := fmt.Sprintf("%s/%d", qr.conf.Url, table.Number)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		qr.logger.Errorw("table token: new request creation error",
			logx.LogField{Key: "table", Value: table.Number},
			logx.LogField{Key: "err", Value: err})
		return err
	}
	req.Header.Add("TableToken", table.Token)
	resp, err := qr.client.Do(req)
	fmt.Println(resp)
	if err != nil {
		qr.logger.Errorw("table token: request error",
			logx.LogField{Key: "table", Value: table.Number},
			logx.LogField{Key: "err", Value: err})
		return err
	}
	fmt.Println(resp.Status)
	if resp.Status != "200 OK" {
		qr.logger.Errorw("table token: validation error",
			logx.LogField{Key: "table", Value: table.Number},
			logx.LogField{Key: "token", Value: table.Token})
		return fmt.Errorf("table token: validation error")
	}
	return nil
}
