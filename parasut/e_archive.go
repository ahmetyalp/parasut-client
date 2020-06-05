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
	Note                   string            `jsonapi:"attr,note"`
	ExciseDutyCodes        []ExciseDutyCodes `jsonapi:"attr,excise_duty_codes,omitempty"`
	SalesInvoice           *SalesInvoice     `jsonapi:"relation,sales_invoice"`
}

type ExciseDutyCodes struct {
	Product             int64  `json:"product,omitempty"`
	SalesExciseDutyCode string `json:"sales_excise_duty_code,omitempty"`
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

	r, err := req.Get(BASE_URL+"v4/"+e_archive.client.CompanyID+"/e_archives/"+id, header, params)

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

	r, err := req.Post(BASE_URL+"v4/"+e_archive.client.CompanyID+"/e_archives", header, req.BodyJSON(body))

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
