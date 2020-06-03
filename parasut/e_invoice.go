package parasut

import (
	"errors"
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
		log.Fatal(err)
		return nil, err
	}

	if statusCode := r.Response().StatusCode; statusCode == 401 || statusCode == 403 {
		return nil, errors.New("unauthorized")
	}

	result := new(EInvoice)

	err = jsonapi.UnmarshalPayload(r.Response().Body, result)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return result, nil
}

func (e_invoice *EInvoice) iamdocument() {

}
