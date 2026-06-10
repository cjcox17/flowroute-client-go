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

	fmt.Println("Flowroute Send SMS Example")
	fmt.Println("===========================")
	fmt.Println()

	// Send a simple SMS
	fmt.Println("Sending SMS...")
	resp, err := client.Messages().Send(ctx).
		To("+19511231234").
		From("+18441231234").
		Body("Hello from the Flowroute Go SDK!").
		Send()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Message ID: %s\n", resp.Data.ID)
	fmt.Printf("Type: %s\n", resp.Data.Type)

	if resp.Data.Attributes.PriceDetails != nil {
		pd := resp.Data.Attributes.PriceDetails
		fmt.Printf("Base Rate:    %s\n", pd.BaseRate)
		fmt.Printf("Charged Cost: %s\n", pd.ChargedCost)
		fmt.Printf("Segments:     %d\n", pd.SegmentCount)
		fmt.Printf("Surcharge:    %s\n", pd.SurchargeRate)
	}

	fmt.Printf("\nLinks:\n")
	fmt.Printf("  Self: %s\n", resp.Data.Links.Self)
}
