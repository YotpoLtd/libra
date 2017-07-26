package nomad

import (
	"errors"
	"os"
	"strconv"

	api "github.com/hashicorp/nomad/api"
)

// NewClient will create a instance of a nomad API Client
func NewClient(c Config) (*api.Client, error) {
	nomadDefaultConfig := api.DefaultConfig()
	envAddress := os.Getenv("NOMAD_ADDRESS")

	if envAddress != "" {
		nomadDefaultConfig.Address = envAddress
	} else {
		nomadDefaultConfig.Address = c.Address
	}

	client, err := api.NewClient(nomadDefaultConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Scale increases or decreases the count of a task group
func Scale(client *api.Client, jobID, groupID string, scale, min, max int) (string, int, error) {
	job, _, _ := client.Jobs().Info(jobID, &api.QueryOptions{})
	oldCount := *job.TaskGroups[0].Count
	newCount := oldCount + scale
	if newCount < min || newCount > max {
		return "", oldCount, errors.New("the new group count (" + strconv.Itoa(newCount) + ") is outside of the configured range (" + strconv.Itoa(min) + "-" + strconv.Itoa(max) + ")")
	}
	job.TaskGroups[0].Count = &newCount
	resp, _, _ := client.Jobs().Register(job, &api.WriteOptions{})
	return resp.EvalID, newCount, nil
}

// Restart restarts a job to get the latest docker image
func Restart(client *api.Client, jobID string) (string, error) {
	evalID, _, err := client.Jobs().ForceEvaluate(jobID, &api.WriteOptions{})
	if err != nil {
		return "", err
	}
	return evalID, nil
}

// SetCapacity sets the count of a task group
func SetCapacity(client *api.Client, jobID, groupID string, count, min, max int) (string, int, error) {
	job, _, _ := client.Jobs().Info(jobID, &api.QueryOptions{})
	oldCount := *job.TaskGroups[0].Count
	if count < min || count > max {
		return "", oldCount, errors.New("the desired count (" + strconv.Itoa(count) + ") is outside of the configured range (" + strconv.Itoa(min) + "-" + strconv.Itoa(max) + ")")
	}
	job.TaskGroups[0].Count = &count
	resp, _, _ := client.Jobs().Register(job, &api.WriteOptions{})
	return resp.EvalID, count, nil
}
