package main

import (
	"fmt"
	"log"
	"os"
	"parasut-client/parasut"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	client := parasut.Client{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		UserName:     os.Getenv("USERNAME"),
		Password:     os.Getenv("PASSWORD"),
		CompanyID:    "31169",
		AutoRefresh:  true,
	}

	client.Connect()

	inboxes, error := client.EInvoiceInbox().All("7171717171")

	if error != nil {
		log.Fatal(error)
		return
	}

	fmt.Println(inboxes[0].Vkn)

	time.Sleep(500 * time.Second)
}
