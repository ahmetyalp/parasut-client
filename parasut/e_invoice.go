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

type EInvoice struct {
	client       *Client
	ID           string `jsonapi:"primary,e_invoices"`
	ExternalID   string `jsonapi:"attr,external_id"`
	UUID         string `jsonapi:"attr,uuid"`
	EnvUUID      string `jsonapi:"attr,env_uuid"`
	FromAddress  string `jsonapi:"attr,from_address"`
	FromVkn      string `jsonapi:"attr,from_vkn"`
	ToAddress    string `jsonapi:"attr,to_address"`
	ToVkn        string `jsonapi:"attr,to_vkn"`
	Direction    string `jsonapi:"attr,direction"`
	Note         string `jsonapi:"attr,note"`
	ResponseType string `jsonapi:"attr,response_type"`
	ContactName  string `jsonapi:"attr,contact_name"`
	Scenario     string `jsonapi:"attr,scenario"`
	Status       string `jsonapi:"attr,status"`
	IssueDate    string `jsonapi:"attr,isseu_date"`
	IsExpired    bool   `jsonapi:"attr,is_expired"`
	IsAnswerable bool   `jsonapi:"attr,is_answerable"`
	NetTolal     int64  `jsonapi:"attr,net_tolal"`
	Currency     string `jsonapi:"attr,currency"`
	ItemType     string `jsonapi:"attr,item_type"`
	SalesInvoice *SalesInvoice
	PurchaseBill *PurchaseBill
}

type EInvoiceParams struct {
	ID                     string `jsonapi:"primary,e_invoices"`
	VatWithholdingCode     string `jsonapi:"attr,vat_withholding_code,omitempty"`
	VatExemptionReasonCode string `jsonapi:"attr,vat_exemption_reason_code,omitempty"`
	VatExemptionReason     string `jsonapi:"attr,vat_exemption_reason,omitempty"`
	Note                   string `jsonapi:"attr,note,omitempty"`
	// use ExciseDutyCodes from e_archive.go
	ExciseDutyCodes []ExciseDutyCodes `jsonapi:"attr,excise_duty_codes,omitempty"`
	Scenario        string            `jsonapi:"attr,scenario"`
	To              string            `jsonapi:"attr,to"`
	SalesInvoice    *SalesInvoice     `jsonapi:"relation,invoice"`
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

	r, err := req.Get(e_invoice.urlBuilder(id), header, params)

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

	result := new(EInvoice)

	err = jsonapi.UnmarshalPayload(body, result)

	// assign invoice
	if strings.Contains(params["include"].(string), "invoice") {
		// parse as regular json with jsonapi structs
		var raw jsonapi.OnePayload
		json.Unmarshal(bodyBytes, &raw)

		invoiceRelation := raw.Data.Relationships["invoice"].(map[string]interface{})["data"]
		if invoiceRelation == nil {
			return result, nil
		}

		// get type and id of invoice
		invoiceType := invoiceRelation.(map[string]interface{})["type"]
		invoiceID := invoiceRelation.(map[string]interface{})["id"]

		// loop for included entities, find invoice
		for _, val := range raw.Included {
			if val.ID == invoiceID && val.Type == invoiceType {
				bytearr, _ := json.Marshal(jsonapi.OnePayload{Data: val})
				body = ioutil.NopCloser(bytes.NewBuffer(bytearr))
				switch invoiceType {
				case "sales_invoices":
					result.SalesInvoice = new(SalesInvoice)
					jsonapi.UnmarshalPayload(body, result.SalesInvoice)
				case "purchase_bills":
					result.PurchaseBill = new(PurchaseBill)
					jsonapi.UnmarshalPayload(body, result.PurchaseBill)
				}
				break
			}
		}
	}

	return result, nil
}

func (e_invoice *EInvoice) New(params *EInvoiceParams) (*TrackableJob, error) {
	header := req.Header{
		"Authorization": "Bearer " + e_invoice.client.AccessToken,
	}

	body, _ := jsonapi.Marshal(params)
	body.(*jsonapi.OnePayload).Included = nil

	r, err := req.Post(e_invoice.urlBuilder(), header, req.BodyJSON(body))

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

	err = jsonapi.UnmarshalPayload(r.Response().Body, result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	result.client = e_invoice.client
	return result, nil
}

func (e_invoice *EInvoice) Pdf(id string) (*EDocumentPdfResponse, error) {
	header := req.Header{
		"Authorization": "Bearer " + e_invoice.client.AccessToken,
	}

	r, err := req.Get(e_invoice.urlBuilder(id, "pdf"), header)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = HandleHTTPStatus(r.Response())

	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := new(EDocumentPdfResponse)

	err = jsonapi.UnmarshalPayload(r.Response().Body, result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}

func (e_invoice *EInvoice) urlBuilder(params ...string) string {
	return e_invoice.client.UrlBuilder(append([]string{"e_invoices"}, params...)...)
}
