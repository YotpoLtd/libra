package api

import (
	"net/http"
	"os"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
	"github.com/underarmour/libra/config"
	"github.com/underarmour/libra/nomad"
)

type RestartRequest struct {
	Job string `json:"job"`
}

type RestartResponse struct {
	Eval string `json:"eval"`
}

func NewRestartRequest(job string) *RestartRequest {
	return &RestartRequest{
		Job: job,
	}
}

func RestartHandler(w rest.ResponseWriter, r *rest.Request) {
	var t RestartRequest
	err := r.DecodeJsonPayload(&t)
	if err != nil {
		log.Errorln(err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	config, err := config.NewConfig(os.Getenv("LIBRA_CONFIG_DIR"))
	if err != nil {
		log.Errorf("Failed to read or parse config file: %s", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Info("Loaded and parsed configuration file")
	n, err := nomad.NewClient(config.Nomad)
	if err != nil {
		log.Errorf("Failed to create Nomad Client: %s", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Info("Successfully created Nomad Client")

	evalID, err := nomad.Restart(n, t.Job)
	if err != nil {
		log.Error("Problem restarting the job " + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		log.Infoln("Restarted it! Evaluation " + evalID)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		respBody := &RestartResponse{
			Eval: evalID,
		}

		w.WriteJson(respBody)
	}
}
