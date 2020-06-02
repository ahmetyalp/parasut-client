package parasut

import (
	"errors"
	"log"
	"reflect"

	"github.com/google/jsonapi"
	"github.com/imroc/req"
)

type EInvoiceInbox struct {
	client *Client
	ID     int    `jsonapi:"primary,e_invoice_inboxes"`
	Vkn    string `jsonapi:"attr,vkn"`
}

func (c *Client) EInvoiceInbox() *EInvoiceInbox {
	return &EInvoiceInbox{
		client: c,
	}
}

func (e_invoice_inbox *EInvoiceInbox) All(vkn string) ([]*EInvoiceInbox, error) {
	params := req.QueryParam{
		"filter[vkn]": vkn,
	}

	header := req.Header{
		"Authorization": "Bearer " + e_invoice_inbox.client.AccessToken,
	}
	var r *req.Resp
	var err error
	r, err = req.Get(BASE_URL+"v4/"+e_invoice_inbox.client.CompanyID+"/e_invoice_inboxes", params, header)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if statusCode := r.Response().StatusCode; statusCode == 401 || statusCode == 403 {
		return nil, errors.New("unauthorized")
	}
	var data []interface{}
	data, err = jsonapi.UnmarshalManyPayload(r.Response().Body, reflect.TypeOf(e_invoice_inbox))

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	result := make([]*EInvoiceInbox, len(data))

	for i, val := range data {
		result[i] = val.(*EInvoiceInbox)
	}

	return result, nil
}
