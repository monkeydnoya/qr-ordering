package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"qr-ordering-service/internal/integration/auth"
	"qr-ordering-service/internal/types"
)

type MyResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (mrw *MyResponseWriter) Write(p []byte) (int, error) {
	return mrw.buf.Write(p)
}

type Middleware struct {
	QrAuthProvider auth.QrAuth
}

func NewMiddleware(qrAuthProvider auth.QrAuth) *Middleware {
	middleware := Middleware{QrAuthProvider: qrAuthProvider}
	return &middleware
}

func (mw *Middleware) Handler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order types.Order
		var table types.Table
		body, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(body, &order)
		if err != nil {
			http.Error(w, "Error when parsing body", http.StatusUnauthorized)
			return
		}

		table.Number = order.Table
		table.Token = r.Header.Get("TableToken")

		if table.Token == "" {
			http.Error(w, "TableToken header missing", http.StatusUnauthorized)
			return
		}
		if err := mw.QrAuthProvider.ValidateTable(table); err != nil {
			http.Error(w, "Token validation error", http.StatusUnauthorized)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	}
}
