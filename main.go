package main

import (
	"fmt"
	"os"
	"parasut-client/parasut"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	client := parasut.Client{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		UserName:     os.Getenv("USERNAME"),
		Password:     os.Getenv("PASSWORD"),
		CompanyID:    "31169",
		AutoRefresh:  true,
	}

	err := client.Connect()

	if err != nil {
		fmt.Println(err)
		return
	}

	// inboxes, _ := client.EInvoiceInbox().All("7171717171")

	// fmt.Printf("%d\n", len(inboxes))
	// fmt.Printf("%#v\n", inboxes[0])
	// fmt.Printf("%#v\n", inboxes[1])

	// eInvoice, _ := client.EInvoice().Find("53456")

	// fmt.Printf("%#v\n", eInvoice)

	// salesInvoice, _ := client.SalesInvoice().Find("513907", "active_e_document", "details", "contact")

	// fmt.Printf("%#v\n", salesInvoice)
	// fmt.Printf("%#v\n", salesInvoice.Contact)
	// fmt.Printf("%#v\n", salesInvoice.ActiveEInvoice)
	// fmt.Printf("%#v\n", salesInvoice.ActiveEArchive)
	// fmt.Printf("%#v\n", len(salesInvoice.SalesInvoiceDetails))
	// fmt.Printf("%#v\n", salesInvoice.SalesInvoiceDetails[0])

	contact := parasut.Contact{
		ID: "401080",
	}

	product := parasut.Product{
		ID: "548459",
	}

	details := parasut.SalesInvoiceDetails{
		Product:  &product,
		Quantity: "6",
	}
	sales_invoice := parasut.SalesInvoice{
		IssueDate:           "2020-01-01",
		DueDate:             "2020-01-01",
		Contact:             &contact,
		WithholdingRate:     "0",
		ItemType:            "invoice",
		Currency:            "TRL",
		ExchangeRate:        "1.0",
		InvoiceDiscountType: "amount",
		InvoiceDiscount:     "0",
		SalesInvoiceDetails: []*parasut.SalesInvoiceDetails{&details},
	}

	client.SalesInvoice().New(&sales_invoice)

	// contact, _ := client.Contact().Find("1982805")

	// fmt.Printf("%#v\n", contact)

	// eArchive, _ := client.EArchive().Find("155730", "sales_invoice")

	// fmt.Printf("%#v\n\n", eArchive)
	// fmt.Printf("%#v\n", eArchive.SalesInvoice)

	// sales_invoice := &parasut.SalesInvoice{ID: "2495610"}
	// params := parasut.EArchiveParams{
	// 	Note:            "ali",
	// 	SalesInvoice:    sales_invoice,
	// 	ExciseDutyCodes: []parasut.ExciseDutyCodes{{Product: 2}},
	// 	Shipment:        parasut.Shipment{Name: "yurtici", Date: "2015-01-01"},
	// }

	// tj, err := client.EArchive().New(&params)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Printf("%#v\n\n", tj)

	// for tj.Status == "pending" || tj.Status == "running" {
	// 	tj.Poll()
	// 	fmt.Printf("%#v\n\n", tj)
	// 	time.Sleep(5 * time.Second)
	// }

	// pdf, err := client.EArchive().Pdf("155730")

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Printf("%#v\n\n", pdf)

	// e_invoice, err := client.EInvoice().Find("53450", "invoice")

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Printf("%#v\n", e_invoice)
	// fmt.Printf("%#v\n", e_invoice.SalesInvoice)

	// params := parasut.EInvoiceParams{
	// 	VatWithholdingCode: "23",
	// 	SalesInvoice:       sales_invoice,
	// 	ExciseDutyCodes:    []parasut.ExciseDutyCodes{{Product: 2}},
	// 	Scenario:           "basic",
	// 	To:                 "ali",
	// }

	// tj, err := client.EInvoice().New(&params)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Printf("%#v\n\n", tj)

	// for tj.Status == "pending" || tj.Status == "running" {
	// 	tj.Poll()
	// 	fmt.Printf("%#v\n\n", tj)
	// 	time.Sleep(5 * time.Second)
	// }

	// pdf, err := client.EInvoice().Pdf("53456")

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Printf("%#v\n\n", pdf)

}
