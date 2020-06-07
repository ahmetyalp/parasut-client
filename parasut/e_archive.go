package parasut

import (
	"log"
	"strings"

	"github.com/google/jsonapi"
	"github.com/imroc/req"
)

type EArchive struct {
	client           *Client
	ID               string        `jsonapi:"primary,e_archives"`
	InvoiceNumber    string        `jsonapi:"attr,invoice_number"`
	UUID             string        `jsonapi:"attr,uuid"`
	Vkn              string        `jsonapi:"attr,vkn"`
	Note             string        `jsonapi:"attr,note"`
	IsPrinted        bool          `jsonapi:"attr,is_printed"`
	Status           string        `jsonapi:"attr,status"`
	PrintedAt        string        `jsonapi:"attr,printed_at"`
	CancellableUntil string        `jsonapi:"attr,cancellable_until"`
	IsSigned         bool          `jsonapi:"attr,is_signed"`
	SalesInvoice     *SalesInvoice `jsonapi:"relation,sales_invoice"`
}

type EArchiveParams struct {
	ID                     string            `jsonapi:"primary,e_archives"`
	VatWithholdingCode     string            `jsonapi:"attr,vat_withholding_code,omitempty"`
	VatExemptionReasonCode string            `jsonapi:"attr,vat_exemption_reason_code,omitempty"`
	VatExemptionReason     string            `jsonapi:"attr,vat_exemption_reason,omitempty"`
	Note                   string            `jsonapi:"attr,note,omitempty"`
	ExciseDutyCodes        []ExciseDutyCodes `jsonapi:"attr,excise_duty_codes,omitempty"`
	InternetSale           InternetSale      `jsonapi:"attr,internet_sale,omitempty"`
	Shipment               Shipment          `jsonapi:"attr,shipment,omitempty"`
	SalesInvoice           *SalesInvoice     `jsonapi:"relation,sales_invoice"`
}

type ExciseDutyCodes struct {
	Product             int64  `json:"product,omitempty"`
	SalesExciseDutyCode string `json:"sales_excise_duty_code,omitempty"`
}

type InternetSale struct {
	Url             string `json:"url,omitempty"`
	PaymentType     string `json:"payment_type,omitempty"`
	PaymentPlatform string `json:"payment_platform,omitempty"`
	PaymentDate     string `json:"payment_date,omitempty"`
}

type Shipment struct {
	Title  string `json:"title,omitempty"`
	Vkn    string `json:"vkn,omitempty"`
	Name   string `json:"name,omitempty"`
	Tcknow string `json:"tckn,omitempty"`
	Date   string `json:"date,omitempty"`
}

type EDocumentPdfResponse struct {
	ID        string `jsonapi:"primary,e_document_pdfs"`
	Url       string `jsonapi:"attr,url"`
	ExpiresAt string `jsonapi:"attr,expires_at"`
}

func (c *Client) EArchive() *EArchive {
	return &EArchive{
		client: c,
	}
}

func (e_archive *EArchive) Find(id string, include ...string) (*EArchive, error) {
	params := req.QueryParam{
		"include": strings.Join(include, ","),
	}

	header := req.Header{
		"Authorization": "Bearer " + e_archive.client.AccessToken,
	}

	r, err := req.Get(e_archive.urlBuilder(id), header, params)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = HandleHTTPStatus(r.Response())

	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := new(EArchive)

	err = jsonapi.UnmarshalPayload(r.Response().Body, result)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}

func (e_archive *EArchive) New(params *EArchiveParams) (*TrackableJob, error) {
	header := req.Header{
		"Authorization": "Bearer " + e_archive.client.AccessToken,
	}

	body, _ := jsonapi.Marshal(params)
	body.(*jsonapi.OnePayload).Included = nil

	r, err := req.Post(e_archive.urlBuilder(), header, req.BodyJSON(body))

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
	result.client = e_archive.client
	return result, nil
}

func (e_archive *EArchive) Pdf(id string) (*EDocumentPdfResponse, error) {
	header := req.Header{
		"Authorization": "Bearer " + e_archive.client.AccessToken,
	}

	r, err := req.Get(e_archive.urlBuilder(id, "pdf"), header)

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

func (e_archive *EArchive) urlBuilder(params ...string) string {
	return e_archive.client.UrlBuilder(append([]string{"e_archives"}, params...)...)
}
