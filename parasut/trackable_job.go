package parasut

import (
	"log"

	"github.com/google/jsonapi"
	"github.com/imroc/req"
)

type TrackableJob struct {
	client *Client
	ID     string   `jsonapi:"primary,trackable_jobs"`
	Status string   `jsonapi:"attr,status"`
	Errors []string `jsonapi:"attr,errors"`
}

func (c *Client) TrackableJob() *TrackableJob {
	return &TrackableJob{
		client: c,
	}
}

func (tj *TrackableJob) PollWithId(id string) (*TrackableJob, error) {
	header := req.Header{
		"Authorization": "Bearer " + tj.client.AccessToken,
	}

	r, err := req.Get(BASE_URL+"v4/"+tj.client.CompanyID+"/trackable_jobs/"+id, header)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = HandleHTTPStatus(r.Response())

	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := new(TrackableJob)

	client := tj.client
	err = jsonapi.UnmarshalPayload(r.Response().Body, result)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	tj.client = client
	return result, nil
}

func (tj *TrackableJob) Poll() error {
	header := req.Header{
		"Authorization": "Bearer " + tj.client.AccessToken,
	}

	r, err := req.Get(BASE_URL+"v4/"+tj.client.CompanyID+"/trackable_jobs/"+tj.ID, header)

	if err != nil {
		log.Println(err)
		return err
	}

	err = HandleHTTPStatus(r.Response())

	if err != nil {
		log.Println(err)
		return err
	}

	client := tj.client
	err = jsonapi.UnmarshalPayload(r.Response().Body, tj)

	if err != nil {
		log.Println(err)
		return err
	}

	tj.client = client
	return nil
}
