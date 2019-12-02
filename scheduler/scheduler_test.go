package scheduler

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
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
			Schedule: "16 16 1 * *",
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
	name := uuid.New().String()
	req := &UpsertJobRequest{
		ProjectID:   "gcpug-ds2bq-dev",
		Location:    "asia-northeast1",
		Name:        name,
		Description: "unit test sample",
		Schedule:    "16 16 1 * *",
		TimeZone:    "Asia/Tokyo",
		Target: &JobHttpTarget{
			Uri:                "https://gcpug-ds2bq-tf572eohna-an.a.run.app/api/v1/datastore-export/",
			HttpMethod:         scheduler.HttpMethod_GET,
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
	_, err = c.Upsert(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	{
		req.Schedule = "20 20 1 * *"
		job, err := c.Upsert(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if e, g := req.Schedule, job.Schedule; e != g {
			t.Errorf("Schedule want %v got %v", e, g)
		}
	}
	{
		req.TimeZone = "UTC"
		job, err := c.Upsert(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if e, g := req.TimeZone, job.TimeZone; e != g {
			t.Errorf("TimeZone want %v got %v", e, g)
		}
	}
	{
		req.Description = "hello scheduler"
		job, err := c.Upsert(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if e, g := req.Description, job.Description; e != g {
			t.Errorf("Description want %v got %v", e, g)
		}
	}
	{
		req.Target.Uri = "https://gcpug-ds2bq-tf572eohna-an.a.run.app/api/v1/datastore-export/?hoge=fuga"
		job, err := c.Upsert(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if e, g := req.Target.Uri, job.GetHttpTarget().Uri; e != g {
			t.Errorf("Target.Uri want %v got %v", e, g)
		}
		if e, g := req.Target.Uri, job.GetHttpTarget().GetOidcToken().Audience; e != g {
			t.Errorf("Target.OidcToken.Audience want %v got %v", e, g)
		}
	}
	{
		req.Target.HttpMethod = scheduler.HttpMethod_POST
		job, err := c.Upsert(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if e, g := req.Target.HttpMethod, job.GetHttpTarget().HttpMethod; e != g {
			t.Errorf("Target.HttpMethod want %v got %v", e, g)
		}
	}
	{
		m := make(map[string]string)
		m["X-Hoge"] = "hoge"
		req.Target.Headers = m
		job, err := c.Upsert(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		m["User-Agent"] = "Google-Cloud-Scheduler" // defaultで入るものを比較する時に合わせるために追加
		if e, g := req.Target.Headers, job.GetHttpTarget().Headers; !reflect.DeepEqual(e, g) {
			t.Errorf("Target.Headers want %v got %v", e, g)
		}
	}
	{
		req.Target.Body = []byte(`{"hoge":"fuga"}`)
		job, err := c.Upsert(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if e, g := req.Target.Body, job.GetHttpTarget().GetBody(); !reflect.DeepEqual(e, g) {
			t.Errorf("Target.Body want %v got %v", string(e), string(g))
		}
	}
	{
		req.Target.OidcServiceAccount = "scheduler2@gcpug-ds2bq-dev.iam.gserviceaccount.com"
		job, err := c.Upsert(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if e, g := req.Target.OidcServiceAccount, job.GetHttpTarget().GetOidcToken().ServiceAccountEmail; !reflect.DeepEqual(e, g) {
			t.Errorf("Target.OidcServiceAccount want %v got %v", string(e), string(g))
		}
	}
}
