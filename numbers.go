package flowroute

// NumbersClient provides access to phone number related APIs
type NumbersClient struct {
	client *Client
}

// Purchaseable returns a client for searching and purchasing phone numbers
func (n *NumbersClient) Purchaseable() *PurchaseableNumbersClient {
	return &PurchaseableNumbersClient{client: n.client}
}

// PhoneNumbers returns a client for managing account phone numbers
func (n *NumbersClient) PhoneNumbers() *PhoneNumbersClient {
	return &PhoneNumbersClient{client: n.client}
}
