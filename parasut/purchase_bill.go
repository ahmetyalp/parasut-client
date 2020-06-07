package parasut

type PurchaseBill struct {
	client *Client
	ID     string `jsonapi:"primary,purchase_bill"`
}

func (c *Client) PurchaseBill() *PurchaseBill {
	return &PurchaseBill{
		client: c,
	}
}
