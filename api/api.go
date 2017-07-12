package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	rootcerts "github.com/hashicorp/go-rootcerts"
)

// Much of this is cribbed from https://github.com/hashicorp/nomad/blob/master/api/api.go

// QueryOptions are used to parameterize a query
type QueryOptions struct {
	// WaitIndex is used to enable a blocking query. Waits
	// until the timeout or the next index is reached
	WaitIndex uint64

	// WaitTime is used to bound the duration of a wait.
	// Defaults to that of the Config, but can be overridden.
	WaitTime time.Duration

	// Set HTTP parameters on the query.
	Params map[string]string
}

// WriteOptions are used to parameterize a write
type WriteOptions struct {
	// Providing a datacenter overwrites the region provided
	// by the Config
	Region string
}

// HTTPBasicAuth is used to authenticate http client with HTTP Basic Authentication
type HTTPBasicAuth struct {
	// Username to use for HTTP Basic Authentication
	Username string

	// Password to use for HTTP Basic Authentication
	Password string
}

// Config is used to configure the creation of a client
type Config struct {
	// Address is the address of the Libra agent
	Address string

	// Region to use. If not provided, the default agent region is used.
	Region string

	// HTTPClient is the client to use. Default will be
	// used if not provided.
	HTTPClient *http.Client

	// HTTPAuth is the auth info to use for http access.
	HTTPAuth *HTTPBasicAuth

	// WaitTime limits how long a Watch will block. If not provided,
	// the agent default values will be used.
	WaitTime time.Duration

	// TLSConfig provides the various TLS related configurations for the http
	// client
	TLSConfig *TLSConfig
}

// TLSConfig contains the parameters needed to configure TLS on the HTTP client
// used to communicate with Libra.
type TLSConfig struct {
	// CACert is the path to a PEM-encoded CA cert file to use to verify the
	// Libra server SSL certificate.
	CACert string

	// CAPath is the path to a directory of PEM-encoded CA cert files to verify
	// the Libra server SSL certificate.
	CAPath string

	// ClientCert is the path to the certificate for Libra communication
	ClientCert string

	// ClientKey is the path to the private key for Libra communication
	ClientKey string

	// TLSServerName, if set, is used to set the SNI host when connecting via
	// TLS.
	TLSServerName string

	// Insecure enables or disables SSL verification
	Insecure bool
}

// DefaultConfig returns a default configuration for the client
func DefaultConfig() *Config {
	config := &Config{
		Address:    "http://127.0.0.1:8080",
		HTTPClient: cleanhttp.DefaultClient(),
		TLSConfig:  &TLSConfig{},
	}
	transport := config.HTTPClient.Transport.(*http.Transport)
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	if addr := os.Getenv("LIBRA_ADDR"); addr != "" {
		config.Address = addr
	}
	if auth := os.Getenv("LIBRA_HTTP_AUTH"); auth != "" {
		var username, password string
		if strings.Contains(auth, ":") {
			split := strings.SplitN(auth, ":", 2)
			username = split[0]
			password = split[1]
		} else {
			username = auth
		}

		config.HTTPAuth = &HTTPBasicAuth{
			Username: username,
			Password: password,
		}
	}

	// Read TLS specific env vars
	if v := os.Getenv("LIBRA_CACERT"); v != "" {
		config.TLSConfig.CACert = v
	}
	if v := os.Getenv("LIBRA_CAPATH"); v != "" {
		config.TLSConfig.CAPath = v
	}
	if v := os.Getenv("LIBRA_CLIENT_CERT"); v != "" {
		config.TLSConfig.ClientCert = v
	}
	if v := os.Getenv("LIBRA_CLIENT_KEY"); v != "" {
		config.TLSConfig.ClientKey = v
	}
	if v := os.Getenv("LIBRA_SKIP_VERIFY"); v != "" {
		if insecure, err := strconv.ParseBool(v); err == nil {
			config.TLSConfig.Insecure = insecure
		}
	}

	return config
}

// ConfigureTLS applies a set of TLS configurations to the the HTTP client.
func (c *Config) ConfigureTLS() error {
	if c.HTTPClient == nil {
		return fmt.Errorf("config HTTP Client must be set")
	}

	var clientCert tls.Certificate
	foundClientCert := false
	if c.TLSConfig.ClientCert != "" || c.TLSConfig.ClientKey != "" {
		if c.TLSConfig.ClientCert != "" && c.TLSConfig.ClientKey != "" {
			var err error
			clientCert, err = tls.LoadX509KeyPair(c.TLSConfig.ClientCert, c.TLSConfig.ClientKey)
			if err != nil {
				return err
			}
			foundClientCert = true
		} else {
			return fmt.Errorf("Both client cert and client key must be provided")
		}
	}

	clientTLSConfig := c.HTTPClient.Transport.(*http.Transport).TLSClientConfig
	rootConfig := &rootcerts.Config{
		CAFile: c.TLSConfig.CACert,
		CAPath: c.TLSConfig.CAPath,
	}
	if err := rootcerts.ConfigureTLS(clientTLSConfig, rootConfig); err != nil {
		return err
	}

	clientTLSConfig.InsecureSkipVerify = c.TLSConfig.Insecure

	if foundClientCert {
		clientTLSConfig.Certificates = []tls.Certificate{clientCert}
	}
	if c.TLSConfig.TLSServerName != "" {
		clientTLSConfig.ServerName = c.TLSConfig.TLSServerName
	}

	return nil
}

// Client provides a client to the Libra API
type Client struct {
	config Config
}

// NewClient returns a new client
func NewClient(config *Config) (*Client, error) {
	// bootstrap the config
	defConfig := DefaultConfig()

	if config.Address == "" {
		config.Address = defConfig.Address
	} else if _, err := url.Parse(config.Address); err != nil {
		return nil, fmt.Errorf("invalid address '%s': %v", config.Address, err)
	}

	if config.HTTPClient == nil {
		config.HTTPClient = defConfig.HTTPClient
	}

	client := &Client{
		config: *config,
	}
	return client, nil
}

// NewRequest sends a request with the desired parameters to the server
func (c *Client) NewRequest(path, method string, body interface{}) (*http.Response, error) {
	switch method {
	case "post":
		b, err := encodeBody(body)
		if err != nil {
			return nil, err
		}
		resp, err := http.Post(c.config.Address+path, "application/json; charset=utf-8", b)
		if err != nil {
			return nil, err
		}
		return resp, nil
	case "get":
		resp, err := http.Get(c.config.Address + path)
		if err != nil {
			return nil, err
		}
		return resp, nil
	default:
		return nil, errors.New("invalid method type " + method)
	}

}

// encodeBody is used to encode a request body
func encodeBody(obj interface{}) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf, nil
}
