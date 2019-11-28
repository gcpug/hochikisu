package scheduler

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/genproto/googleapis/cloud/scheduler/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestClient_Update(t *testing.T) {
	ctx := context.Background()
	c, err := NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	job, err := c.Update(ctx, &scheduler.UpdateJobRequest{
		Job: &scheduler.Job{
			Name:        "projects/gcpug-ds2bq-dev/locations/asia-northeast1/jobs/hoge",
			Description: "",
			Target: &scheduler.Job_HttpTarget{
				HttpTarget: &scheduler.HttpTarget{
					Uri:        "https://{YOUR_DS2BQ_CLOUD_RUN_URI}/api/v1/datastore-export/",
					HttpMethod: scheduler.HttpMethod_POST,
					Headers:    nil,
					Body:       nil,
					AuthorizationHeader: &scheduler.HttpTarget_OidcToken{
						OidcToken: &scheduler.OidcToken{
							ServiceAccountEmail: "scheduler@$DS2BQ_PROJECT_ID.iam.gserviceaccount.com",
							Audience:            "https://{YOUR_DS2BQ_CLOUD_RUN_URI}/api/v1/datastore-export/",
						},
					},
				},
			},
			Schedule: "16 16 * * *",
			TimeZone: "Asia/Tokyo",
		},
		UpdateMask: nil,
	})
	if err != nil {
		t.Logf("%T", err)
		if status.Code(err) != codes.NotFound {
			t.Fatal(err)
		}
	}
	fmt.Println(job)
}

func TestClient_Upsert(t *testing.T) {
	req := &UpsertJobRequest{
		ProjectID:   "gcpug-ds2bq-dev",
		Location:    "asia-northeast1",
		Name:        "hoge",
		Description: "unit test sample",
		Schedule:    "16 16 * * *",
		TimeZone:    "Asia/Tokyo",
		Target: &JobHttpTarget{
			Uri:                "https://gcpug-ds2bq-tf572eohna-an.a.run.app/api/v1/datastore-export/",
			HttpMethod:         scheduler.HttpMethod_POST,
			Headers:            nil,
			Body:               nil,
			OidcServiceAccount: "scheduler@gcpug-ds2bq-dev.iam.gserviceaccount.com",
		},
	}

	ctx := context.Background()
	c, err := NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	job, err := c.Upsert(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", job)
}
