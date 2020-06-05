package parasut

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/google/jsonapi"
	"github.com/imroc/req"
)

type SalesInvoice struct {
	client         *Client
	ID             string   `jsonapi:"primary,sales_invoices"`
	NetTotal       string   `jsonapi:"attr,net_total"`
	Contact        *Contact `jsonapi:"relation,contact"`
	ActiveEInvoice *EInvoice
	ActiveEArchive *EArchive
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
		log.Println(err)
		return nil, err
	}

	err = HandleHTTPStatus(r.Response())

	if err != nil {
		log.Println(err)
		return nil, err
	}

	// copy body
	body := r.Response().Body
	bodyBytes, _ := ioutil.ReadAll(body)
	body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	result := new(SalesInvoice)
	err = jsonapi.UnmarshalPayload(body, result)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	// assign active e-invoice/archive
	if strings.Contains(params["include"].(string), "active_e_document") {
		// parse as regular json with jsonapi structs
		var raw jsonapi.OnePayload
		json.Unmarshal(bodyBytes, &raw)

		// get type and id of active e-document
		activeEDocumentType := raw.Data.Relationships["active_e_document"].(map[string]interface{})["data"].(map[string]interface{})["type"]
		activeEDocumentID := raw.Data.Relationships["active_e_document"].(map[string]interface{})["data"].(map[string]interface{})["id"]

		// loop for included entities, find active e-document
		for _, val := range raw.Included {
			if val.ID == activeEDocumentID && val.Type == activeEDocumentType {
				bytearr, _ := json.Marshal(jsonapi.OnePayload{Data: val})
				body = ioutil.NopCloser(bytes.NewBuffer(bytearr))
				switch activeEDocumentType {
				case "e_invoices":
					result.ActiveEInvoice = new(EInvoice)
					jsonapi.UnmarshalPayload(body, result.ActiveEInvoice)
				case "e_archives":
					result.ActiveEArchive = new(EArchive)
					jsonapi.UnmarshalPayload(body, result.ActiveEArchive)
				}
				break
			}
		}
	}

	return result, nil
}
