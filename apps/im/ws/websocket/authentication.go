package websocket

// 鉴权

import (
	"fmt"
	"net/http"
	"time"
)

type Authentication interface {
	Auth(w http.ResponseWriter, r *http.Request) bool
	UserId(r *http.Request) string
}

type authentication struct{}

func (*authentication) Auth(w http.ResponseWriter, r *http.Request) bool {
	return true
}
func (*authentication) UserId(r *http.Request) string {
	query := r.URL.Query()
	if query != nil && query["userId"] != nil {
		return fmt.Sprintf("%v", query["userId"])
	}

	return fmt.Sprintf("%v", time.Now().UnixMilli())
}
