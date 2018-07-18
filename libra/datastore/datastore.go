package datastore

import (
	"github.com/YotpoLtd/libra/config"
	nomad "github.com/hashicorp/nomad/api"
)

type Store interface {
	Get(jobName string) (nomad.Job, error)
	Put(jobName string, jobDescription nomad.Job) error
}

func InitializeStore(config *config.RootConfig) error {

}
