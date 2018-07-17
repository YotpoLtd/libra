package command

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"time"

	"github.com/YotpoLtd/libra/api"
	"github.com/YotpoLtd/libra/backend"
	"github.com/YotpoLtd/libra/config"
	"github.com/YotpoLtd/libra/nomad"
	"github.com/YotpoLtd/libra/structs"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/mitchellh/cli"
	"github.com/sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"
)

type libra struct {
	configDir string
	config    *config.RootConfig
}

// ServerCommand is a Command implementation prints the version.
type ServerCommand struct {
	libra *libra
	Ui    cli.Ui
}

func (c *ServerCommand) Help() string {
	helpText := `
Usage: libra server [options]
  Run a Libra server. The other commands require a server to be configured.
`
	return strings.TrimSpace(helpText)
}

func Setup() (*libra, error) {
	lc := &libra{
		configDir: "/etc/libra",
	}

	if configDir := os.Getenv("LIBRA_CONFIG_DIR"); configDir != "" {
		lc.configDir = configDir
	}

	info, err := os.Stat(lc.configDir)
	if os.IsNotExist(err) {
		logrus.Errorf("Configuration Directory %s does not exists", lc.configDir)
		return nil, err
	}

	if !info.IsDir() {
		logrus.Errorf("Specified Configuration Directory %s is not a directory", lc.configDir)
	}

	lc.config, err = config.ParseConfig(lc.configDir)
	if err != nil {
		logrus.Errorf("Failed to read or parse config file: %s", err)
	}

	return lc, nil
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

func (c *ServerCommand) Run(args []string) int {

	libra, err := Setup()
	if err != nil {
		logrus.Errorf("Boo-boo")
	}

	s := rest.NewApi()
	logger := logrus.New()
	w := logger.Writer()
	defer w.Close()

	loggingMw := &rest.AccessLogApacheMiddleware{
		Logger: log.New(w, "[access] ", 0),
	}

	mw := []rest.Middleware{
		loggingMw,
		&rest.ContentTypeCheckerMiddleware{},
		&rest.GzipMiddleware{},
		&rest.JsonIndentMiddleware{},
		&rest.PoweredByMiddleware{},
		&rest.RecorderMiddleware{},
		&rest.RecoverMiddleware{
			EnableResponseStackTrace: true,
		},
		&rest.TimerMiddleware{},
	}

	s.Use(mw...)
	router, err := rest.MakeRouter(
		rest.Post("/scale", api.ScaleHandler),
		rest.Post("/capacity", api.CapacityHandler),
		rest.Post("/grafana", api.GrafanaHandler),
		rest.Get("/backends", api.BackendsHandler),
		rest.Get("/ping", api.PingHandler),
		rest.Get("/", api.HomeHandler),
		rest.Post("/restart", api.RestartHandler),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	s.SetApp(router)

	cr, _, err := libra.loadRules()
	if err != nil {
		logrus.Errorf("Problem with the Libra server: %s", err)
		return 1
	}
	cr.Start()

	err = http.ListenAndServe(":8646", s.MakeHandler())
	if err != nil {
		logrus.Errorf("Problem with the Libra server: %s", err)
		return 1
	}
	return 0
}

func (c *ServerCommand) Synopsis() string {
	return "Run a Libra server"
}
