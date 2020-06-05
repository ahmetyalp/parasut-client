package main

import (
	"fmt"
	"os"
	"parasut-client/parasut"
	"time"

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

	inboxes, _ := client.EInvoiceInbox().All("7171717171")

	fmt.Printf("%d\n", len(inboxes))
	fmt.Printf("%#v\n", inboxes[0])
	fmt.Printf("%#v\n", inboxes[1])

	eInvoice, _ := client.EInvoice().Find("53456")

	fmt.Printf("%#v\n", eInvoice)

	salesInvoice, _ := client.SalesInvoice().Find("2495627", "active_e_document")

	fmt.Printf("%#v\n", salesInvoice)
	fmt.Printf("%#v\n", salesInvoice.Contact)
	fmt.Printf("%#v\n", salesInvoice.ActiveEInvoice)
	fmt.Printf("%#v\n", salesInvoice.ActiveEArchive)

	contact, _ := client.Contact().Find("1982805")

	fmt.Printf("%#v\n", contact)

	eArchive, _ := client.EArchive().Find("155746", "sales_invoice")

	fmt.Printf("%#v\n\n", eArchive)
	fmt.Printf("%#v\n", eArchive.SalesInvoice)

	sales_invoice := &parasut.SalesInvoice{ID: "2495627"}
	params := parasut.EArchiveParams{
		Note:            "ali",
		SalesInvoice:    sales_invoice,
		ExciseDutyCodes: []parasut.ExciseDutyCodes{parasut.ExciseDutyCodes{Product: 2}},
	}

	tj, err := client.EArchive().New(&params)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%#v\n\n", tj)

	for tj.Status == "pending" || tj.Status == "running" {
		tj.Poll()
		fmt.Printf("%#v\n\n", tj)
		time.Sleep(5 * time.Second)
	}
}
