package main

import (
	"context"
	"fmt"
	"log"
	"os"

	flowroute "github.com/cjcox17/flowroute-client-go"
)

func main() {
	apiKey := os.Getenv("FLOWROUTE_API_KEY")
	if apiKey == "" {
		log.Fatal("FLOWROUTE_API_KEY environment variable not set")
	}

	apiSecret := os.Getenv("FLOWROUTE_API_SECRET")
	if apiSecret == "" {
		log.Fatal("FLOWROUTE_API_SECRET environment variable not set")
	}

	client := flowroute.NewClient(apiKey, apiSecret)
	ctx := context.Background()

	fmt.Println("Flowroute Send MMS Example")
	fmt.Println("===========================")
	fmt.Println()

	// Send an MMS with text and a media attachment
	fmt.Println("Sending MMS...")
	resp, err := client.Messages().Send(ctx).
		To("+12061231234").
		From("+18441231234").
		Body("Check out this image!").
		IsMMS(true).
		MediaURLs("https://sbc.example.com/invoice.png").
		Send()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Message ID: %s\n", resp.Data.ID)
	fmt.Printf("Type: %s\n", resp.Data.Type)

	if resp.Data.Attributes.PriceDetails != nil {
		pd := resp.Data.Attributes.PriceDetails
		fmt.Printf("Charged Cost: %s\n", pd.ChargedCost)
	}

	fmt.Printf("\nLinks:\n")
	fmt.Printf("  Self: %s\n", resp.Data.Links.Self)
}
