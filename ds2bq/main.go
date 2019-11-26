package ds2bq

import (
	"context"

	"gopkg.in/yaml.v2"
)

type SchedulerJob struct {
	Schedule                string  `json:"schedule"`
	URI                     string  `json:"uri"`
	Body                    *Config `json:"messageBody"`
	OIDCServiceAccountEmail string  `json:"oidcServiceAccountEmail" yaml:"oidcServiceAccountEmail"`
}

type Config struct {
	ProjectID         string `json:"projectID" yaml:"projectID"`
	OutputGCSFilePath string `json:"outputGCSFilePath" yaml:"outputGCSFilePath"`
	AllKinds          bool   `json:"allKinds" yaml:"allKinds"`
	BQLoadProjectID   string `json:"bqLoadProjectID" yaml:"bqLoadProjectID"`
	BQLoadDatasetID   string `json:"bqLoadDatasetID" yaml:"bqLoadDatasetID"`
}

func ParseYaml(ctx context.Context, data []byte) ([]*SchedulerJob, error) {
	var jobs []*SchedulerJob
	if err := yaml.Unmarshal(data, &jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}
