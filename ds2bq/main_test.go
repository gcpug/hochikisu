package ds2bq

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestOutputYaml(t *testing.T) {
	var jobs []*SchedulerJob
	job := &SchedulerJob{
		Schedule: "16 16 * * *",
		URI:      "https://{YOUR_DS2BQ_CLOUD_RUN_URI}/api/v1/datastore-export/",
		Body: &Config{
			ProjectID:         "datastore-project",
			OutputGCSFilePath: "gs://datastore-project-ds2bq-test",
			AllKinds:          true,
			BQLoadProjectID:   "datastore-project",
			BQLoadDatasetID:   "ds2bq_test",
		},
		OIDCServiceAccountEmail: "scheduler@$DS2BQ_PROJECT_ID.iam.gserviceaccount.com",
	}
	jobs = append(jobs, job)
	o, err := yaml.Marshal(jobs)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(o))
}

func TestParseYaml(t *testing.T) {
	fn := filepath.Join("testdata", "yaml.golden")
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Fatal(err)
	}

	jobs, err := ParseYaml(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}
	if e, g := 1, len(jobs); e != g {
		t.Errorf("jobs.length want %v got %v", e, g)
	}

	job := jobs[0]

	if e, g := "16 16 * * *", job.Schedule; e != g {
		t.Errorf("Schedule want %s got %s", e, g)
	}
	if e, g := "https://{YOUR_DS2BQ_CLOUD_RUN_URI}/api/v1/datastore-export/", job.URI; e != g {
		t.Errorf("URI want %s got %s", e, g)
	}
	if e, g := "scheduler@$DS2BQ_PROJECT_ID.iam.gserviceaccount.com", job.OIDCServiceAccountEmail; e != g {
		t.Errorf("SA want %s got %s", e, g)
	}
	if job.Body == nil {
		t.Fatal("Body is nil")
	}
	if e, g := "datastore-project", job.Body.ProjectID; e != g {
		t.Errorf("Body.ProjectID want %s got %s", e, g)
	}
	if e, g := "gs://datastore-project-ds2bq-test", job.Body.OutputGCSFilePath; e != g {
		t.Errorf("Body.OutputGCSFilePath want %s got %s", e, g)
	}
	if e, g := true, job.Body.AllKinds; e != g {
		t.Errorf("Body.AllKinds want %v got %v", e, g)
	}
	if e, g := "datastore-project", job.Body.BQLoadProjectID; e != g {
		t.Errorf("Body.BQLoadProjectID %s got %s", e, g)
	}
	if e, g := "ds2bq_test", job.Body.BQLoadDatasetID; e != g {
		t.Errorf("Body.BQLoadDataSetID %s got %s", e, g)
	}
}
