package parasut

import (
	"errors"
	"log"

	"github.com/google/jsonapi"
	"github.com/imroc/req"
)

type Contact struct {
	client  *Client
	ID      string `jsonapi:"primary,contacts"`
	Balance string `jsonapi:"attr,balance"`
}

func (c *Client) Contact() *Contact {
	return &Contact{
		client: c,
	}
}

func (contact *Contact) Find(id string) (*Contact, error) {
	header := req.Header{
		"Authorization": "Bearer " + contact.client.AccessToken,
	}

	r, err := req.Get(BASE_URL+"v4/"+contact.client.CompanyID+"/contacts/"+id, header)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if statusCode := r.Response().StatusCode; statusCode == 401 || statusCode == 403 {
		return nil, errors.New("unauthorized")
	}

	result := new(Contact)

	err = jsonapi.UnmarshalPayload(r.Response().Body, result)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return result, nil
}
