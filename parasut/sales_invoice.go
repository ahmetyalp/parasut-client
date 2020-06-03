package parasut

import (
	"errors"
	"log"
	"strings"

	"github.com/google/jsonapi"
	"github.com/imroc/req"
)

type SalesInvoice struct {
	client          *Client
	ID              string    `jsonapi:"primary,sales_invoices"`
	NetTotal        string    `jsonapi:"attr,net_total"`
	Contact         *Contact  `jsonapi:"relation,contact"`
	ActiveEDocument *EInvoice `jsonapi:"relation,active_e_document"`
}

func (c *Client) SalesInvoice() *SalesInvoice {
	return &SalesInvoice{
		client: c,
	}
}

func (sales_invoice *SalesInvoice) Find(id string, include ...string) (*SalesInvoice, error) {
	params := req.QueryParam{
		"include": strings.Join(include, ","),
	}

	header := req.Header{
		"Authorization": "Bearer " + sales_invoice.client.AccessToken,
	}

	r, err := req.Get(BASE_URL+"v4/"+sales_invoice.client.CompanyID+"/sales_invoices/"+id, header, params)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if statusCode := r.Response().StatusCode; statusCode == 401 || statusCode == 403 {
		return nil, errors.New("unauthorized")
	}

	result := new(SalesInvoice)

	err = jsonapi.UnmarshalPayload(r.Response().Body, result)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return result, nil
}
