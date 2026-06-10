package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	fmt.Println("Flowroute List Messages Example")
	fmt.Println("================================")
	fmt.Println()

	// Look up messages from the last 7 days
	endDate := time.Now().UTC()
	startDate := endDate.Add(-7 * 24 * time.Hour)

	fmt.Printf("Searching messages from %s to %s\n\n",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))

	resp, err := client.Messages().List(ctx).
		StartDate(startDate).
		EndDate(endDate).
		Limit(10).
		Send()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Found %d message(s):\n\n", len(resp.Data))

	for i, msg := range resp.Data {
		attrs := msg.Attributes

		fmt.Printf("%d. %s\n", i+1, msg.ID)
		fmt.Printf("   Direction: %s\n", attrs.Direction)
		fmt.Printf("   From:      %s\n", attrs.From)
		fmt.Printf("   To:        %s\n", attrs.To)
		fmt.Printf("   Status:    %s\n", attrs.Status)
		fmt.Printf("   Body:      %s\n", attrs.Body)
		fmt.Printf("   Timestamp: %s\n", attrs.Timestamp)
		fmt.Printf("   Is MMS:    %t\n", attrs.IsMMS)

		if len(attrs.DeliveryReceipts) > 0 {
			fmt.Printf("   Delivery Receipts:\n")
			for _, dlr := range attrs.DeliveryReceipts {
				fmt.Printf("     Level %d: %s\n", dlr.Level, dlr.Status)
			}
		}

		fmt.Println()
	}

	if resp.Links.Next != nil {
		fmt.Println("More results available. Use Offset() to paginate.")
	}
}
