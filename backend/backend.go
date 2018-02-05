package backend

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/underarmour/libra/config"
	"github.com/underarmour/libra/structs"
)

// ConfiguredBackends struct
type ConfiguredBackends map[string]structs.Backender

// InitializeBackends loads valid backends into a map
func InitializeBackends(backends map[string]structs.Backend) (ConfiguredBackends, error) {
	configuredBackends := make(ConfiguredBackends, 0)

	for name, backend := range backends {
		backendType := backend.Kind

		if backendType == "" {
			return nil, fmt.Errorf("Missing backend type for '%s'", name)
		}

		switch backendType {
		case "cloudwatch":
			c, err := config.NewConfig(os.Getenv("LIBRA_CONFIG_DIR"))
			if err != nil {
				log.Errorf("Failed to read or parse config file: %s", err)
				return nil, err
			}

			conf := c.Backends[name]

			connection, err := NewCloudWatchBackend(name, CloudWatchConfig{
				Kind:   conf.Kind,
				Name:   conf.Name,
				Region: conf.Region,
			})
			if err != nil {
				return nil, fmt.Errorf("Bad configuration for %s: %s", name, err)
			}

			configuredBackends[name] = connection

		case "graphite":
			c, err := config.NewConfig(os.Getenv("LIBRA_CONFIG_DIR"))
			if err != nil {
				log.Errorf("Failed to read or parse config file: %s", err)
				return nil, err
			}

			conf := c.Backends[name]

			password := conf.Password
			if password == "" {
				password = os.Getenv("GRAPHITE_PASSWORD")
			}
			connection, err := NewGraphiteBackend(name, GraphiteConfig{
				Kind:     conf.Kind,
				Name:     conf.Name,
				Host:     conf.Host,
				Username: conf.Username,
				Password: password,
			})
			if err != nil {
				return nil, fmt.Errorf("Bad configuration for %s: %s", name, err)
			}

			configuredBackends[name] = connection

		case "influxdb":
			c, err := config.NewConfig(os.Getenv("LIBRA_CONFIG_DIR"))
			if err != nil {
				log.Errorf("Failed to read or parse config file: %s", err)
				return nil, err
			}

			conf := c.Backends[name]

			username := conf.Password
			if username == "" {
				username = os.Getenv("INFLUX_USERNAME")
			}

			password := conf.Password
			if password == "" {
				password = os.Getenv("INFLUX_PASSWORD")
			}
			connection, err := NewInfluxDbBackend(name, InfluxDbConfig{
				Kind:     conf.Kind,
				Name:     conf.Name,
				Addr:     conf.Addr,
				Timeout:  conf.Timeout,
				Username: username,
				Password: password,
			})
			if err != nil {
				return nil, fmt.Errorf("Bad configuration for %s: %s", name, err)
			}

			configuredBackends[name] = connection

		default:
			log.Fatalf("unknown backend type '%s' for backend %s", backendType, name)
			return nil, fmt.Errorf("unknown backend %s", backendType)
		}
	}

	return configuredBackends, nil
}
