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

	fmt.Println("Flowroute Simple Search Example")
	fmt.Println("================================")
	fmt.Println()

	// Search for phone numbers in the 206 area code
	fmt.Println("Searching for numbers in area code 206...")
	response, err := client.Numbers().Purchaseable().
		SearchAvailable(ctx).
		StartsWith("206").
		State("WA").
		NumberType("longcode").
		OrderBy("monthly_cost", "asc").
		Limit(10).
		Send()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("\nFound %d available phone numbers:\n\n", len(response.Data))

	for i, number := range response.Data {
		state := "N/A"
		if number.Attributes.State != nil {
			state = *number.Attributes.State
		}

		rateCenter := "N/A"
		if number.Attributes.RateCenter != nil {
			rateCenter = *number.Attributes.RateCenter
		}

		fmt.Printf("%d. %s\n", i+1, number.Attributes.Value)
		fmt.Printf("   Monthly Cost: $%.2f\n", number.Attributes.MonthlyCost)
		fmt.Printf("   Setup Cost:   $%.2f\n", number.Attributes.SetupCost)
		fmt.Printf("   State:        %s\n", state)
		fmt.Printf("   Rate Center:  %s\n\n", rateCenter)
	}

	if response.Links.Next != nil {
		fmt.Println("More results available. Use Offset() to paginate.")
	}
}
