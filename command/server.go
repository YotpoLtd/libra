package command

import (
	"log"
	"net/http"
	"strings"

	"github.com/YotpoLtd/libra/api"
	"github.com/YotpoLtd/libra/libra"

	"flag"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/mitchellh/cli"
	"github.com/sirupsen/logrus"
)

// ServerCommand is a Command implementation prints the version.
type ServerCommand struct {
	Ui cli.Ui
}

func (c *ServerCommand) Help() string {
	helpText := `
Usage: libra server [options]
  Run a Libra server. The other commands require a server to be configured.
`
	return strings.TrimSpace(helpText)
}

func (c *ServerCommand) Run(args []string) int {
	serverFlags := flag.NewFlagSet("server", flag.ContinueOnError)
	confPath := "/etc/libra"
	serverFlags.StringVar(&confPath, "config", confPath, "Config directory for Libra")
	serverFlags.Parse(args)

	libra, err := libra.Setup(confPath)
	if err != nil {
		logrus.Errorf("Server start failed =(")
		return 1
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
		rest.Get("/ping", api.PingHandler),
		rest.Get("/", api.HomeHandler),
		rest.Post("/restart", api.RestartHandler),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	s.SetApp(router)

	cr, _, err := libra.LoadRules()
	if err != nil {
		logrus.Errorf("Problem with the Libra server: %s", err)
		return 1
	}
	cr.Start()
	// TODO: Make listener honor ENV/Config
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
