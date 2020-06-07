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
	client                 *Client
	ID                     string   `jsonapi:"primary,sales_invoices"`
	Archived               bool     `jsonapi:"attr,archived"`
	NetTotal               string   `jsonapi:"attr,net_total"`
	GrossTotal             string   `jsonapi:"attr,gross_total"`
	Withholding            string   `jsonapi:"attr,withholding"`
	TotalExciseDuty        string   `jsonapi:"attr,total_excise_duty"`
	TotalCommunicationsTax string   `jsonapi:"attr, total_communications_tax "`
	TotalVat               string   `jsonapi:"attr,total_vat"`
	VatWitholding          string   `jsonapi:"attr,vat_witholding"`
	TotalDiscount          string   `jsonapi:"attr,total_discount"`
	TotalInvoiceDiscount   string   `jsonapi:"attr,total_invoice_discount"`
	BeforeTaxesTotal       string   `jsonapi:"attr,before_taxes_total"`
	Remaining              string   `jsonapi:"attr,remaining"`
	RemainingInTrl         string   `jsonapi:"attr,remaining_in_trl"`
	PaymentStatus          string   `jsonapi:"attr,payment_status"`
	ItemType               string   `jsonapi:"attr,item_type"`
	Description            string   `jsonapi:"attr,description"`
	IssueDate              string   `jsonapi:"attr,issue_date"`
	DueDate                string   `jsonapi:"attr,due_date"`
	InvoiceNo              string   `jsonapi:"attr,invoice_no"`
	InvoiceSeries          string   `jsonapi:"attr,invoice_series"`
	InvoiceID              string   `jsonapi:"attr,invoice_id"`
	Currency               string   `jsonapi:"attr,currency"`
	ExchangeRate           string   `jsonapi:"attr,exchange_rate"`
	WithholdingRate        string   `jsonapi:"attr,withholding_rate"`
	VatWitholdingRate      string   `jsonapi:"attr,vat_witholding_rate"`
	InvoiceDiscountType    string   `jsonapi:"attr,invoice_discount_tyoe"`
	InvoiceDiscount        string   `jsonapi:"attr,invoice_discount"`
	BillingAddress         string   `jsonapi:"attr,billing_address"`
	BillingPhone           string   `jsonapi:"attr,billing_phone"`
	BillingFax             string   `jsonapi:"attr,billing_fax"`
	TaxOffice              string   `jsonapi:"attr,tax_office"`
	TaxNumber              string   `jsonapi:"attr,tax_number"`
	City                   string   `jsonapi:"attr,city"`
	District               string   `jsonapi:"attr,district"`
	IsAbroad               bool     `jsonapi:"attr,is_abroad"`
	OrderNo                string   `jsonapi:"attr,order_no"`
	OrderDate              string   `jsonapi:"attr,order_date"`
	ShipmentAddress        string   `jsonapi:"attr,shipment_address"`
	ShipmentIncluded       bool     `jsonapi:"attr,shipment_included"`
	Contact                *Contact `jsonapi:"relation,contact"`
	SalesInvoiceDetails    []*SalesInvoiceDetails
	ActiveEInvoice         *EInvoice
	ActiveEArchive         *EArchive
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

	r, err := req.Get(sales_invoice.urlBuilder(id), header, params)

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
		activeEDocumentRelation := raw.Data.Relationships["active_e_document"].(map[string]interface{})["data"]
		if activeEDocumentRelation == nil {
			return result, nil
		}

		activeEDocumentType := activeEDocumentRelation.(map[string]interface{})["type"]
		activeEDocumentID := activeEDocumentRelation.(map[string]interface{})["id"]

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

func (sales_invoice *SalesInvoice) New(sales_invoice_params *SalesInvoice) (*SalesInvoice, error) {
	header := req.Header{
		"Authorization": "Bearer " + sales_invoice.client.AccessToken,
	}

	// set body
	body, _ := jsonapi.Marshal(sales_invoice_params)
	body.(*jsonapi.OnePayload).Included = nil
	// details should be under the relationships, not included
	details, _ := jsonapi.Marshal(sales_invoice_params.SalesInvoiceDetails)
	body.(*jsonapi.OnePayload).Data.Relationships["details"] = details

	r, err := req.Post(sales_invoice.urlBuilder(), header, req.BodyJSON(body))

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = HandleHTTPStatus(r.Response())

	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := new(SalesInvoice)
	err = jsonapi.UnmarshalPayload(r.Response().Body, result)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func (sales_invoice *SalesInvoice) urlBuilder(params ...string) string {
	return sales_invoice.client.UrlBuilder(append([]string{"sales_invoices"}, params...)...)
}
