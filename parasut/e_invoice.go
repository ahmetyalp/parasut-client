package parasut

import (
	"log"
	"strings"

	"github.com/google/jsonapi"
	"github.com/imroc/req"
)

type EInvoice struct {
	client      *Client
	ID          string `jsonapi:"primary,e_invoices"`
	ExternalID  string `jsonapi:"attr,external_id"`
	UUID        string `jsonapi:"attr,uuid"`
	EnvUUID     string `jsonapi:"attr,env_uuid"`
	FromAddress string `jsonapi:"attr,from_address"`
	FromVkn     string `jsonapi:"attr,from_vkn"`
}

func (c *Client) EInvoice() *EInvoice {
	return &EInvoice{
		client: c,
	}
}

func (e_invoice *EInvoice) Find(id string, include ...string) (*EInvoice, error) {
	params := req.QueryParam{
		"include": strings.Join(include, ","),
	}

	header := req.Header{
		"Authorization": "Bearer " + e_invoice.client.AccessToken,
	}

	r, err := req.Get(BASE_URL+"v4/"+e_invoice.client.CompanyID+"/e_invoices/"+id, header, params)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = HandleHTTPStatus(r.Response())

	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := new(EInvoice)

	err = jsonapi.UnmarshalPayload(r.Response().Body, result)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}
