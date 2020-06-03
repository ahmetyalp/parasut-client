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

	client.Connect()

	// inboxes, _ := client.EInvoiceInbox().All("7171717171")

	// fmt.Printf("%+v\n", inboxes[0])

	// eInvoice, _ := client.EInvoice().Find("53456")

	// fmt.Printf("%+v\n", eInvoice)

	salesInvoice, _ := client.SalesInvoice().Find("2495610", "contact", "active_e_document")

	fmt.Printf("%+v\n", salesInvoice)
	fmt.Printf("%+v\n", *salesInvoice.Contact)
	fmt.Printf("%+v\n", salesInvoice.ActiveEDocument)

	// contact, _ := client.Contact().Find("1982805")

	// fmt.Printf("%+v\n", contact)
}
