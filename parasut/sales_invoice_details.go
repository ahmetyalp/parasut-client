package parasut

type SalesInvoiceDetails struct {
	ID                    string   `jsonapi:"primary,sales_invoice_details"`
	Quantity              string   `jsonapi:"attr,quantity,omitempty"`
	UnitPrice             string   `jsonapi:"attr,unit_price,omitempty"`
	VatRate               string   `jsonapi:"attr,vat_rate,omitempty"`
	DiscountType          string   `jsonapi:"attr,discount_type,omitempty"`
	DiscountValue         string   `jsonapi:"attr,discount_value,omitempty"`
	ExciseDutyType        string   `jsonapi:"attr,excise_duty_type,omitempty"`
	ExciseDutyValue       string   `jsonapi:"attr,excise_duty_value,omitempty"`
	CommunicationsTaxRate string   `jsonapi:"attr,communications_tax_rate,omitempty"`
	Description           string   `jsonapi:"attr,description,omitempty"`
	Product               *Product `jsonapi:"relation,product"`
}
