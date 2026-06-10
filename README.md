# Flowroute Client for Go

A modern, type-safe Go client library for the [Flowroute API](https://www.flowroute.com/).

[![Go Reference](https://pkg.go.dev/badge/github.com/cjcox17/flowroute-client-go.svg)](https://pkg.go.dev/github.com/cjcox17/flowroute-client-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/cjcox17/flowroute-client-go)](https://goreportcard.com/report/github.com/cjcox17/flowroute-client-go)
[![License](https://img.shields.io/badge/license-MIT%2FApache--2.0-blue.svg)](LICENSE-MIT)

## Features

- 🎯 Type-safe API interactions with comprehensive error handling
- ⚡ Context support for cancellation and timeouts
- 🔧 Builder pattern for complex request construction
- 📦 Zero external dependencies (uses only Go standard library)
- 🔐 Secure Basic Authentication
- 📖 Extensive documentation and examples
- ✅ Idiomatic Go design

## Weirdness

- I've managed to find a delete endpoint to remove numbers, but it's not in the navigation. [https://developer.flowroute.com/api/numbers/v2.0/release-phone-number-from-account/]
- `SearchAvailable()` can have a null State or Rate Center, but this isn't documented.

## Installation

```bash
go get github.com/cjcox17/flowroute-client-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"
	
	flowroute "github.com/cjcox17/flowroute-client-go"
)

func main() {
	client := flowroute.NewClient("your-api-key", "your-api-secret")
	ctx := context.Background()

	response, err := client.Numbers().Purchaseable().
		ListAvailableAreaCodes(ctx).
		Limit(10).
		Send()
	if err != nil {
		log.Fatal(err)
	}

	for _, areaCode := range response.Data {
		fmt.Println("Area Code:", areaCode.ID)
	}
}
```

## Examples

### List Available Area Codes

```go
ctx := context.Background()

response, err := client.Numbers().Purchaseable().
	ListAvailableAreaCodes(ctx).
	Limit(10).
	Offset(0).
	MaxSetupCost(5.0).
	Send()
if err != nil {
	log.Fatal(err)
}

for _, ac := range response.Data {
	fmt.Printf("Area Code: %s\n", ac.ID)
}
```

### List Available Exchanges

```go
response, err := client.Numbers().Purchaseable().
	ListAvailableExchanges(ctx).
	AreaCode(415).
	Limit(10).
	Offset(0).
	MaxSetupCost(3.0).
	Send()
if err != nil {
	log.Fatal(err)
}
```

### List Available Tiers

```go
response, err := client.Numbers().Purchaseable().
	ListAvailableTiers(ctx).
	Limit(10).
	Offset(0).
	Send()
if err != nil {
	log.Fatal(err)
}
```

### Search for Purchasable Phone Numbers

```go
response, err := client.Numbers().Purchaseable().
	SearchAvailable(ctx).
	StartsWith("206").
	Contains("555").
	EndsWith("1234").
	State("WA").
	RateCenter("seattle").
	NumberType("longcode").
	ISOCountry("US").
	OrderBy("monthly_cost", "asc").
	Limit(10).
	Offset(0).
	Send()
if err != nil {
	log.Fatal(err)
}

for _, number := range response.Data {
	fmt.Printf("Number: %s - $%.2f/month\n", 
		number.Attributes.Value, 
		number.Attributes.MonthlyCost)
}
```

### Purchase a Phone Number

```go
response, err := client.Numbers().PhoneNumbers().
	Purchase(ctx, "12065551234").
	CNAMLookupsEnabled(true).
	MessagingEnabled(true).
	MessagingCallbackURL("https://example.com/callback").
	Send()
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Purchased: %s\n", response.Data.Attributes.Value)
```

### List Account Phone Numbers

```go
response, err := client.Numbers().PhoneNumbers().
	List(ctx).
	StartsWith("206").
	Limit(10).
	Send()
if err != nil {
	log.Fatal(err)
}

for _, number := range response.Data {
	fmt.Printf("Number: %s (%s)\n", 
		number.Attributes.Value, 
		number.Attributes.NumberType)
}
```

### Delete a Number

```go
err := client.Numbers().PhoneNumbers().
	RemovePhoneNumber(ctx, "12065551234")
if err != nil {
	log.Fatal(err)
}
fmt.Println("Number removed successfully")
```

### Send an SMS

```go
resp, err := client.Messages().Send(ctx).
	To("+19511231234").
	From("+18441231234").
	Body("Hello from Flowroute!").
	Send()
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Message ID: %s\n", resp.Data.ID)
fmt.Printf("Cost: %s\n", resp.Data.Attributes.PriceDetails.ChargedCost)
```

### Send an MMS

```go
resp, err := client.Messages().Send(ctx).
	To("+12061231234").
	From("+18441231234").
	Body("Check out this image").
	IsMMS(true).
	MediaURLs("https://example.com/photo.png").
	Send()
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Message ID: %s\n", resp.Data.ID)
```

### Look Up a Message Detail Record

```go
mdr, err := client.Messages().Get(ctx, "mdr2-abc123").Send()
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Status: %s\n", mdr.Attributes.Status)
fmt.Printf("From: %s\n", mdr.Attributes.From)
fmt.Printf("To: %s\n", mdr.Attributes.To)
fmt.Printf("Body: %s\n", mdr.Attributes.Body)
```

### List Messages

```go
import "time"

start := time.Now().Add(-7 * 24 * time.Hour)
end := time.Now()

resp, err := client.Messages().List(ctx).
	StartDate(start).
	EndDate(end).
	Limit(10).
	Send()
if err != nil {
	log.Fatal(err)
}

for _, msg := range resp.Data {
	fmt.Printf("%s: %s -> %s\n", msg.ID, msg.Attributes.From, msg.Attributes.To)
}
```

### Set Account Callback URL

```go
_, err := client.Messages().SetCallback(ctx, flowroute.SMSCallback).
	URL("https://example.com/sms-webhook").
	Send()
if err != nil {
	log.Fatal(err)
}
fmt.Println("SMS callback set")
```

### Look Up Account Callback

```go
cb, err := client.Messages().GetCallback(ctx, flowroute.SMSCallback).Send()
if err != nil {
	log.Fatal(err)
}
fmt.Printf("Callback URL: %s\n", cb.Data.Attributes.CallbackURL)
```

### Remove Account Callback

```go
err := client.Messages().RemoveCallback(ctx, flowroute.SMSCallback)
if err != nil {
	log.Fatal(err)
}
fmt.Println("SMS callback removed")
```

### Parse Inbound Message Webhook

```go
// In your HTTP handler where Flowroute POSTs inbound messages
var inbound flowroute.InboundMessage
if err := json.NewDecoder(r.Body).Decode(&inbound); err != nil {
	http.Error(w, err.Error(), http.StatusBadRequest)
	return
}

fmt.Printf("From: %s\n", inbound.Data.Attributes.From)
fmt.Printf("To: %s\n", inbound.Data.Attributes.To)
fmt.Printf("Body: %s\n", inbound.Data.Attributes.Body)
fmt.Printf("Is MMS: %t\n", inbound.Data.Attributes.IsMMS)

// Access MMS media attachments
for _, media := range inbound.Data.Included {
	fmt.Printf("Media: %s (%s)\n", media.Attributes.FileName, media.Attributes.URL)
}
```

### Parse DLR Webhook

```go
// In your HTTP handler where Flowroute POSTs delivery receipts
var dlr flowroute.DeliveryReceiptPayload
if err := json.NewDecoder(r.Body).Decode(&dlr); err != nil {
	http.Error(w, err.Error(), http.StatusBadRequest)
	return
}

fmt.Printf("Message ID: %s\n", dlr.Data.ID)
fmt.Printf("Status: %s\n", dlr.Data.Attributes.Status)
fmt.Printf("Level: %d\n", dlr.Data.Attributes.Level)
```

### Error Handling

```go
import (
	"errors"
	flowroute "github.com/cjcox17/flowroute-client-go"
)

err := client.Numbers().PhoneNumbers().
	RemovePhoneNumber(ctx, "12065551234")
if err != nil {
	var apiErr *flowroute.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("API Error: Status %d\n", apiErr.StatusCode)
		if apiErr.Detail != nil {
			fmt.Printf("Detail: %s\n", apiErr.Detail.Error())
		}
	} else {
		log.Fatal(err)
	}
}
```

### Context with Timeout

```go
import "time"

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

response, err := client.Numbers().Purchaseable().
	ListAvailableAreaCodes(ctx).
	Limit(10).
	Send()
if err != nil {
	log.Fatal(err)
}
```

### Custom HTTP Client

```go
import (
	"net/http"
	"time"
)

httpClient := &http.Client{
	Timeout: 60 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
	},
}

client := flowroute.NewClient("api-key", "api-secret",
	flowroute.WithHTTPClient(httpClient),
)
```

## Running Examples

```bash
export FLOWROUTE_API_KEY="your-api-key"
export FLOWROUTE_API_SECRET="your-api-secret"
```

### Search Examples

```bash
go run examples/simple_search/main.go
go run examples/comprehensive_search/main.go
go run examples/list_account_numbers/main.go
```

### Messaging Examples

```bash
go run examples/send_sms/main.go
go run examples/send_mms/main.go
go run examples/list_messages/main.go
```

### Bulk Purchase Example

Purchase numbers in bulk while excluding numbers listed in a CSV file:

```bash
go run examples/bulk_purchase/main.go \
    -csv excluded_numbers.csv \
    -starts-with 206 \
    -quantity 5 \
    -purchase
```

Flags:
- `-csv` — path to a CSV file of numbers to exclude (plain digits, no header)
- `-starts-with` — search prefix for available numbers
- `-quantity` — how many numbers to purchase
- `-purchase` — **required** to actually buy numbers; without it the run is a dry-run
- `-max-pages` — safety limit for search result pages (default: 50)
- `-page-limit` — results per search request (default: 50)

## Testing

```bash
go test ./...
```

## API Coverage

### Numbers API

- ✅ List available area codes
- ✅ List available exchanges
- ✅ List available tiers
- ✅ Search available phone numbers
- ✅ List account phone numbers
- ✅ Purchase phone number
- ✅ Remove phone number

### Messages API

- ✅ Send SMS
- ✅ Send MMS
- ✅ Look up message detail record (MDR)
- ✅ List messages
- ✅ Set account callback URL
- ✅ Look up account callback
- ✅ Remove account callback
- ✅ Parse inbound message webhook
- ✅ Parse delivery receipt webhook

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

Licensed under either of:

- Apache License, Version 2.0 ([LICENSE-APACHE](LICENSE-APACHE) or http://www.apache.org/licenses/LICENSE-2.0)
- MIT License ([LICENSE-MIT](LICENSE-MIT) or http://opensource.org/licenses/MIT)

at your option.

## Resources

- [Flowroute API Documentation](https://developer.flowroute.com/)
- [Go Package Documentation](https://pkg.go.dev/github.com/cjcox17/flowroute-client-go)
- [Repository](https://github.com/cjcox17/flowroute-client-go)
