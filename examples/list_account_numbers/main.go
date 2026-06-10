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

	fmt.Println("Flowroute Account Phone Numbers")
	fmt.Println("================================")
	fmt.Println()

	// List all phone numbers on the account
	fmt.Println("Fetching account phone numbers...")
	response, err := client.Numbers().PhoneNumbers().
		List(ctx).
		Limit(50).
		Send()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("\nYou have %d phone number(s) on your account:\n\n", len(response.Data))

	for i, number := range response.Data {
		alias := "No alias"
		if number.Attributes.Alias != nil {
			alias = *number.Attributes.Alias
		}

		status := "UNKNOWN"
		if number.Attributes.Status != nil {
			status = *number.Attributes.Status
		}

		messagingEnabled := "No"
		if number.Attributes.MessagingEnabled != nil && *number.Attributes.MessagingEnabled {
			messagingEnabled = "Yes"
		}

		cnamEnabled := "No"
		if number.Attributes.CNAMLookupsEnabled != nil && *number.Attributes.CNAMLookupsEnabled {
			cnamEnabled = "Yes"
		}

		fmt.Printf("%d. %s (%s)\n", i+1, number.Attributes.Value, alias)
		fmt.Printf("   Type:             %s\n", number.Attributes.NumberType)
		fmt.Printf("   Status:           %s\n", status)
		fmt.Printf("   State:            %s\n", number.Attributes.State)
		fmt.Printf("   Rate Center:      %s\n", number.Attributes.RateCenter)
		fmt.Printf("   Messaging:        %s\n", messagingEnabled)
		fmt.Printf("   CNAM Lookups:     %s\n\n", cnamEnabled)
	}

	if response.Links.Next != nil {
		fmt.Println("More results available. Use Offset() to paginate.")
	} else {
		fmt.Println("All phone numbers retrieved.")
	}
}
