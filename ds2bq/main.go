package ds2bq

import (
	"context"
	"encoding/json"

	"github.com/gcpug/hochikisu/scheduler"
	cs "google.golang.org/genproto/googleapis/cloud/scheduler/v1"
	"gopkg.in/yaml.v2"
)

type SchedulerJob struct {
	Name                    string  `json:"name"`
	ProjectID               string  `json:"projectId" yaml:"projectId"`
	Location                string  `json:"location"`
	Description             string  `json:"description"`
	URI                     string  `json:"uri"`
	Schedule                string  `json:"schedule"`
	Timezone                string  `json:"timezone"`
	Body                    *Config `json:"messageBody"`
	OIDCServiceAccountEmail string  `json:"oidcServiceAccountEmail" yaml:"oidcServiceAccountEmail"`
}

type Config struct {
	ProjectID         string   `json:"projectId" yaml:"projectId"`
	AllKinds          bool     `json:"allKinds" yaml:"allKinds"`
	Kinds             []string `json:"kinds"`
	NamespaceIDs      []string `json:"namespaceIds" yaml:"namespaceIds"`
	IgnoreKinds       []string `json:"ignoreKinds" yaml:"ignoreKinds"`
	IgnoreBQLoadKinds []string `json:"ignoreBQLoadKinds" yaml:"ignoreBQLoadKinds"`
	OutputGCSFilePath string   `json:"outputGCSFilePath" yaml:"outputGCSFilePath"`
	BQLoadProjectID   string   `json:"bqLoadProjectId" yaml:"bqLoadProjectId"`
	BQLoadDatasetID   string   `json:"bqLoadDatasetId" yaml:"bqLoadDatasetId"`
}

func ParseYaml(ctx context.Context, data []byte) ([]*SchedulerJob, error) {
	var jobs []*SchedulerJob
	if err := yaml.Unmarshal(data, &jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (job *SchedulerJob) CreateUpsertJobRequest() (*scheduler.UpsertJobRequest, error) {
	body, err := json.Marshal(job.Body)
	if err != nil {
		return nil, err
	}
	return &scheduler.UpsertJobRequest{
		ProjectID:   job.ProjectID,
		Location:    job.Location,
		Name:        job.Name,
		Description: job.Description,
		Schedule:    job.Schedule,
		TimeZone:    job.Timezone,
		Target: &scheduler.JobHttpTarget{
			Uri:                job.URI,
			HttpMethod:         cs.HttpMethod_POST,
			Headers:            nil, // Headerのユースケースがないので、設定していない
			Body:               body,
			OidcServiceAccount: job.OIDCServiceAccountEmail,
		},
	}, nil
}
