package parasut

import (
	"errors"
	"log"
	"strings"

	"github.com/google/jsonapi"
	"github.com/imroc/req"
)

type EArchive struct {
	client        *Client
	ID            string `jsonapi:"primary,e_archives"`
	InvoiceNumber string `jsonapi:"attr,invoice_number"`
	UUID          string `jsonapi:"attr,uuid"`
	Vkn           string `jsonapi:"attr,vkn"`
}

func (c *Client) EArchive() *EArchive {
	return &EArchive{
		client: c,
	}
}

func (e_invoice *EArchive) Find(id string, include ...string) (*EArchive, error) {
	params := req.QueryParam{
		"include": strings.Join(include, ","),
	}

	header := req.Header{
		"Authorization": "Bearer " + e_invoice.client.AccessToken,
	}

	r, err := req.Get(BASE_URL+"v4/"+e_invoice.client.CompanyID+"/e_archives/"+id, header, params)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if statusCode := r.Response().StatusCode; statusCode == 401 || statusCode == 403 {
		return nil, errors.New("unauthorized")
	}

	result := new(EArchive)

	err = jsonapi.UnmarshalPayload(r.Response().Body, result)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return result, nil
}
