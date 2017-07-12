package api

import (
	"github.com/ant0ine/go-json-rest/rest"
)

func HomeHandler(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("welcome to libra, the Nomad auto-scaler")
}
