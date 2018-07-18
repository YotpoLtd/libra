package datastore

import (
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/nomad/api"
	"time"
	"net/http"
)

// consulStore struct description
type consulStore struct {
	client                   *consulapi.Client      // consul client used to interact with Consul
	clientConfig             *consulapi.Config      // consul client config
}

// NewConsulStore will create a new Consul Store Client
func NewConsulStore(name string, config consulConfig) (*consulStore, error) {

}


func (*consulStore) Get(jobName string) (api.Job, error) {
	panic("implement me")
}

func (*consulStore) Put(jobName string, jobDescription api.Job) error {
	panic("implement me")
}
