package api

import (
	"github.com/ant0ine/go-json-rest/rest"
)

func PingHandler(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("pong")
}
