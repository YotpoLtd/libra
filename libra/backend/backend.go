package backend

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/YotpoLtd/libra/structs"
)

// ConfiguredBackends struct
type ConfiguredBackends map[string]structs.Backender

// InitializeBackends loads valid backends into a map
func InitializeBackends(backends map[string]structs.Backend) (ConfiguredBackends, error) {
	configuredBackends := make(ConfiguredBackends, 0)

	for name, backend := range backends {
		backendType := backend.Kind

		if backendType == "" {
			return nil, fmt.Errorf("[ERROR] Missing backend type for '%s'", name)
		}

		switch backendType {
		case "cloudwatch":
			connection, err := NewCloudWatchBackend(name, CloudWatchConfig{
				Kind:   backend.Kind,
				Name:   backend.Name,
				Region: backend.Region,
			})
			if err != nil {
				return nil, fmt.Errorf("[ERROR] Bad configuration for %s: %s", name, err)
			}

			configuredBackends[name] = connection

		case "graphite":

			var graphitePassword string
			if password := os.Getenv("GRAPHITE_PASSWORD"); password != "" {
				graphitePassword = password
			} else {
				graphitePassword = backend.Password
			}

			connection, err := NewGraphiteBackend(name, GraphiteConfig{
				Kind:     backend.Kind,
				Name:     backend.Name,
				Host:     backend.Host,
				Username: backend.Username,
				Password: graphitePassword,
			})
			if err != nil {
				return nil, fmt.Errorf("[ERROR] Bad configuration for %s: %s", name, err)
			}

			configuredBackends[name] = connection

		case "influxdb":
			var influxDbUserName, influxDbPassword string

			if envInfluxDbUserName := os.Getenv("INFLUX_USERNAME"); envInfluxDbUserName != "" {
				influxDbUserName = envInfluxDbUserName
			} else {
				influxDbUserName = backend.Username
			}

			if envInfluxDbPassword := os.Getenv("INFLUX_PASSWORD"); envInfluxDbPassword != "" {
				influxDbPassword = envInfluxDbPassword
			} else {
				influxDbPassword = backend.Password
			}

			userAgent := backend.UserAgent
			if userAgent == "" {
				userAgent = "libra"
			}

			connection, err := NewInfluxDbBackend(name, InfluxDbConfig{
				Addr:      backend.Addr,
				Kind:      backend.Kind,
				Name:      backend.Name,
				Password:  influxDbPassword,
				Timeout:   backend.Timeout,
				UserAgent: userAgent,
				Username:  influxDbUserName,
			})
			if err != nil {
				return nil, fmt.Errorf("[ERROR] Bad configuration for %s: %s", name, err)
			}

			configuredBackends[name] = connection

		default:
			log.Fatalf("unknown backend type '%s' for backend %s", backendType, name)
			return nil, fmt.Errorf("unknown backend %s", backendType)
		}
	}

	return configuredBackends, nil
}
