package parasut

import (
	"log"
	"reflect"

	"github.com/google/jsonapi"
	"github.com/imroc/req"
)

type EInvoiceInbox struct {
	client              *Client
	ID                  string `jsonapi:"primary,e_invoice_inboxes"`
	Vkn                 string `jsonapi:"attr,vkn"`
	EInvoiceAddress     string `jsonapi:"attr, e_invoice_address"`
	Name                string `jsonapi:"attr,name"`
	InboxType           string `jsonapi:"attr,inbox_type"`
	AddressRegisteredAt string `jsonapi:"attr,address_registered_at"`
	RegisteredAt        string `jsonapi:"attr,registered_at"`
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

	r, err := req.Get(BASE_URL+"v4/"+e_invoice_inbox.client.CompanyID+"/e_invoice_inboxes", params, header)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = HandleHTTPStatus(r.Response())

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var data []interface{}
	data, err = jsonapi.UnmarshalManyPayload(r.Response().Body, reflect.TypeOf(e_invoice_inbox))

	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := make([]*EInvoiceInbox, len(data))

	for i, val := range data {
		result[i] = val.(*EInvoiceInbox)
	}

	return result, nil
}
