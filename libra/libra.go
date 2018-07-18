package libra

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/YotpoLtd/libra/libra/backend"
	"github.com/YotpoLtd/libra/libra/config"
	"github.com/YotpoLtd/libra/libra/structs"
	"github.com/hashicorp/nomad/nomad"
	"gopkg.in/robfig/cron.v2"
)

type Libra struct {
	Config *config.Config
}

func Setup(configDir string) (*Libra, error) {
	var newLibra *Libra
	var err error
	newLibra.Config, err = config.ParseConfig(configDir)
	if err != nil {
		logrus.Errorf("Failed to read or parse config file: %s", err)
	}

	return newLibra, nil
}

func (l *libra) loadRules() (*cron.Cron, []cron.EntryID, error) {

	//	n, err := nomad.NewClient(l.config.Nomad)
	//	if err != nil {
	//		log.Fatalf("Failed to create Nomad Client: %s", err)
	//	}
	//	logrus.Info("Successfully created Nomad Client")
	//	dc, err := n.Agent().Datacenter()
	//	if err != nil {
	//		logrus.Fatalf("  Failed to get Nomad DC: %s", err)
	//	}
	//	logrus.Infof("  -> DC: %s", dc)
	backends, err := backend.InitializeBackends(l.config.Backends)
	if err != nil {
		logrus.Fatalf("%s", err)
	}
	logrus.Info("")
	logrus.Infof("Found %d backends", len(backends))
	for name, b := range backends {
		logrus.Infof("  -> %s (%s)", name, b.Info().Kind)
	}
	logrus.Info("")
	logrus.Infof("Found %d jobs", len(l.config.Jobs))

	//err = datastore.InitializeStore(config)
	//if err != nil {
	//	logrus.Fatalf("%s", err)
	//}

	cr := cron.New()
	ids := []cron.EntryID{}

	for _, job := range l.config.Jobs {
		logrus.Infof("  -> Job: %s", job.Name)

		for _, group := range job.Groups {
			logrus.Infof("  --> Group: %s", group.Name)
			logrus.Infof("      min_count = %d", group.MinCount)
			logrus.Infof("      max_count = %d", group.MaxCount)

			for name, rule := range group.Rules {
				cfID, err := cr.AddFunc(rule.Period, createCronFunc(rule, &l.config.Nomad, job.Name, group.Name, group.MinCount, group.MaxCount))
				if err != nil {
					logrus.Errorf("Problem adding autoscaling rule to cron: %s", err)
					return cr, ids, err
				}
				ids = append(ids, cfID)
				logrus.Infof("  ----> Rule: %s", rule.Name)
				if backends[rule.Backend] == nil {
					return cr, ids, fmt.Errorf("Unknown backend: %s (%s)", rule.Backend, name)
				}

				rule.BackendInstance = backends[rule.Backend]
			}
		}
	}
	return cr, ids, nil
}

func createCronFunc(rule *structs.Rule, nomadConf *nomad.Config, job, group string, min, max int) func() {
	return func() {
		n := rand.Intn(10) // offset cron jobs slightly so they don't collide
		time.Sleep(time.Duration(n) * time.Second)
		backend.Work(rule, nomadConf, job, group, min, max)
	}
}
