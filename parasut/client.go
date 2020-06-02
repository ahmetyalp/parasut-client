package parasut

import (
	"log"
	"time"

	"github.com/imroc/req"
)

var BASE_URL string = "http://api.parasut.localhost:3000/"

type Client struct {
	ClientID     string
	ClientSecret string
	UserName     string
	Password     string
	AccessToken  string
	RefreshToken string
	CompanyID    string
	AutoRefresh  bool
}

func (c *Client) Connect() {
	// req.Debug = true
	params := req.QueryParam{
		"client_id":     c.ClientID,
		"client_secret": c.ClientSecret,
		"username":      c.UserName,
		"password":      c.Password,
		"grant_type":    "password",
		"redirect_uri":  "urn:ietf:wg:oauth:2.0:oob",
	}

	r, err := req.Post(BASE_URL+"/oauth/token", params)

	if err != nil {
		log.Fatal(err)
	}
	res := make(map[string]string)
	r.ToJSON(&res)

	c.AccessToken = res["access_token"]
	c.RefreshToken = res["refresh_token"]

	if c.AutoRefresh {
		timer := time.NewTicker(60 * 60 * time.Second)
		go refresher(c, timer)
	}
}

func (c *Client) refresh() {
	params := req.QueryParam{
		"client_id":     c.ClientID,
		"client_secret": c.ClientSecret,
		"grant_type":    "refresh_token",
		"refresh_token": c.RefreshToken,
	}

	r, err := req.Post(BASE_URL+"/oauth/token", params)

	if err != nil {
		log.Fatal(err)
	}
	res := make(map[string]string)
	r.ToJSON(&res)

	c.AccessToken = res["access_token"]
	c.RefreshToken = res["refresh_token"]
}

func refresher(c *Client, timer *time.Ticker) {
	for {
		<-timer.C
		c.refresh()
	}
}
