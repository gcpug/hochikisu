package scheduler

import (
	"context"
	"fmt"

	cs "cloud.google.com/go/scheduler/apiv1"
	"github.com/morikuni/failure"
	"google.golang.org/genproto/googleapis/cloud/scheduler/v1"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInvalidArgument failure.StringCode = "InvalidArgument"

type UpsertJobRequest struct {
	ProjectID   string
	Location    string
	Name        string
	Description string
	Schedule    string
	TimeZone    string
	Target      *JobHttpTarget
}

type JobHttpTarget struct {
	Uri                string
	HttpMethod         scheduler.HttpMethod
	Headers            map[string]string
	Body               []byte
	OidcServiceAccount string
}

func (req *UpsertJobRequest) CreateJobRequest() (*scheduler.CreateJobRequest, error) {
	if req.Target == nil {
		return nil, failure.New(ErrInvalidArgument, failure.Messagef("Target is required"))
	}
	return &scheduler.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", req.ProjectID, req.Location),
		Job: &scheduler.Job{
			Name:        fmt.Sprintf("projects/%s/locations/%s/jobs/%s", req.ProjectID, req.Location, req.Name),
			Description: req.Description,
			Target:      req.Target.JobHttpTarget(),
			Schedule:    req.Schedule,
			TimeZone:    req.TimeZone,
		},
	}, nil
}

func (req *UpsertJobRequest) UpdateJobRequest() (*scheduler.UpdateJobRequest, error) {
	if req.Target == nil {
		return nil, failure.New(ErrInvalidArgument, failure.Messagef("Target is required"))
	}
	return &scheduler.UpdateJobRequest{
		Job: &scheduler.Job{
			Name:        fmt.Sprintf("projects/%s/locations/%s/jobs/%s", req.ProjectID, req.Location, req.Name),
			Description: req.Description,
			Target:      req.Target.JobHttpTarget(),
			Schedule:    req.Schedule,
			TimeZone:    req.TimeZone,
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"description"},
		},
	}, nil
}

func (t *JobHttpTarget) JobHttpTarget() *scheduler.Job_HttpTarget {
	return &scheduler.Job_HttpTarget{
		HttpTarget: &scheduler.HttpTarget{
			Uri:        t.Uri,
			HttpMethod: t.HttpMethod,
			Body:       t.Body,
			AuthorizationHeader: &scheduler.HttpTarget_OidcToken{
				OidcToken: &scheduler.OidcToken{
					ServiceAccountEmail: t.OidcServiceAccount,
					Audience:            t.Uri,
				},
			},
		},
	}
}

type Client struct {
	C *cs.CloudSchedulerClient
}

func NewClient(ctx context.Context) (*Client, error) {
	client, err := cs.NewCloudSchedulerClient(ctx)
	if err != nil {
		return nil, err
	}
	c := Client{
		C: client,
	}
	return &c, nil
}

func (c *Client) Upsert(ctx context.Context, req *UpsertJobRequest) (*scheduler.Job, error) {
	ujReq, err := req.UpdateJobRequest()
	if err != nil {
		return nil, err
	}
	job, err := c.Update(ctx, ujReq)
	if status.Code(err) == codes.NotFound {
		cjReq, err := req.CreateJobRequest()
		if err != nil {
			return nil, err
		}
		job, err := c.Create(ctx, cjReq)
		if err != nil {
			return nil, err
		}
		return job, nil
	} else if err != nil {
		return nil, err
	}
	return job, nil
}

func (c *Client) Create(ctx context.Context, req *scheduler.CreateJobRequest) (*scheduler.Job, error) {
	job, err := c.C.CreateJob(ctx, req)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (c *Client) Update(ctx context.Context, req *scheduler.UpdateJobRequest) (*scheduler.Job, error) {
	job, err := c.C.UpdateJob(ctx, req)
	if err != nil {
		return nil, err
	}
	return job, nil
}
