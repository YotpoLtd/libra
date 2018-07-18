package api

import (
	"github.com/ant0ine/go-json-rest/rest"
)

func HomeHandler(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("Hi, I'm libra, I can scale jobs in Nomad")
}
