package graphite

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"bytes"

	log "github.com/sirupsen/logrus"
)

// Client wraps http.Client so the consumer doesn't have to
type Client struct {
	HTTP     *http.Client
	Host     string
	Username string
	Password string
}

type RenderResponse struct {
	Target     string      `json:"target"`
	Datapoints [][]float64 `json:"datapoints"`
}

// type Datapoint struct {
// 	Value     float64 `json:""`
// 	Timestamp int64   `json:""`
// }

// NewClient creates a new Graphite client, including a custom net/http client
func NewClient(url, username, password string) *Client {
	return &Client{
		HTTP: &http.Client{
			Timeout: time.Second * 10,
		},
		Host:     url,
		Username: username,
		Password: password,
	}
}

// Render makes a call to the Graphite /render endpoint: https://graphite-api.readthedocs.io/en/latest/api.html
func (c *Client) Render(metric string) (RenderResponse, error) {
	var data RenderResponse
	req, err := http.NewRequest("GET", c.Host+"/graphite/render?target="+metric+"&format=json", nil)
	if err != nil {
		log.Errorf("problem creating graphite request: %s", err)
		return data, err
	}
	req.SetBasicAuth(c.Username, c.Password)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		log.Errorf("problem getting graphite response: %s", err)
		return data, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("problem parsing graphite response: %s", err)
		return data, err
	}
	bFlattened := flattenJSON(b)
	json.Unmarshal(bFlattened, &data)
	return data, nil
}

func flattenJSON(b []byte) []byte {
	by := bytes.TrimPrefix(b, []byte("["))
	by = bytes.TrimSuffix(by, []byte("]"))
	return by
}
