package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

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

	fmt.Println("Flowroute API Client - Comprehensive Search Example")
	fmt.Println("====================================================")
	fmt.Println()

	client := flowroute.NewClient(apiKey, apiSecret)
	ctx := context.Background()

	// Step 1: Explore available tiers
	fmt.Println("Step 1: Discovering available pricing tiers...")
	fmt.Println("------------------------------------------------")
	tiersResponse, err := client.Numbers().Purchaseable().
		ListAvailableTiers(ctx).
		Limit(5).
		Send()
	if err != nil {
		log.Fatalf("Error listing tiers: %v", err)
	}

	fmt.Printf("Found %d pricing tiers:\n", len(tiersResponse.Data))
	for _, tier := range tiersResponse.Data {
		fmt.Printf("  • %s\n", tier.ID)
	}
	fmt.Println()

	// Step 2: Find available area codes with affordable setup costs
	fmt.Println("Step 2: Finding affordable area codes (max $3.00 setup cost)...")
	fmt.Println("---------------------------------------------------------------")
	areaCodesResponse, err := client.Numbers().Purchaseable().
		ListAvailableAreaCodes(ctx).
		MaxSetupCost(3.0).
		Limit(5).
		Send()
	if err != nil {
		log.Fatalf("Error listing area codes: %v", err)
	}

	fmt.Printf("Found %d affordable area codes:\n", len(areaCodesResponse.Data))
	for _, areaCode := range areaCodesResponse.Data {
		fmt.Printf("  • Area Code: %s\n", areaCode.ID)
	}
	fmt.Println()

	// Step 3: For each area code, explore available exchanges
	fmt.Println("Step 3: Exploring exchanges in each area code...")
	fmt.Println("------------------------------------------------")

	maxAreaCodes := 3
	if len(areaCodesResponse.Data) < maxAreaCodes {
		maxAreaCodes = len(areaCodesResponse.Data)
	}

	for i := 0; i < maxAreaCodes; i++ {
		areaCode := areaCodesResponse.Data[i]
		acNum, err := strconv.Atoi(areaCode.ID)
		if err != nil {
			continue
		}

		fmt.Printf("\nArea Code: %d\n", acNum)

		exchangesResponse, err := client.Numbers().Purchaseable().
			ListAvailableExchanges(ctx).
			AreaCode(acNum).
			Limit(5).
			Send()
		if err != nil {
			log.Printf("Error listing exchanges for %d: %v", acNum, err)
			continue
		}

		fmt.Printf("  Exchanges available: %d\n", len(exchangesResponse.Data))
		maxExchanges := 3
		if len(exchangesResponse.Data) < maxExchanges {
			maxExchanges = len(exchangesResponse.Data)
		}

		for j := 0; j < maxExchanges; j++ {
			exchange := exchangesResponse.Data[j]
			fmt.Printf("    %d. Exchange %s - %s\n", j+1, exchange.ID, exchange.Type)
		}

		if len(exchangesResponse.Data) > 3 {
			fmt.Printf("    ... and %d more\n", len(exchangesResponse.Data)-3)
		}
	}
	fmt.Println()

	// Step 4: Demonstrate pagination across endpoints
	fmt.Println("Step 4: Demonstrating pagination...")
	fmt.Println("------------------------------------")

	fmt.Println("\n→ Paginating through area codes:")
	page1, err := client.Numbers().Purchaseable().
		ListAvailableAreaCodes(ctx).
		Limit(2).
		Offset(0).
		Send()
	if err != nil {
		log.Fatalf("Error with pagination: %v", err)
	}

	fmt.Printf("  Page 1: %d area codes\n", len(page1.Data))
	for _, ac := range page1.Data {
		fmt.Printf("    - %s\n", ac.ID)
	}

	if page1.Links.Next != nil {
		page2, err := client.Numbers().Purchaseable().
			ListAvailableAreaCodes(ctx).
			Limit(2).
			Offset(2).
			Send()
		if err != nil {
			log.Printf("Error getting page 2: %v", err)
		} else {
			fmt.Printf("  Page 2: %d area codes\n", len(page2.Data))
			for _, ac := range page2.Data {
				fmt.Printf("    - %s\n", ac.ID)
			}
		}
	}
	fmt.Println()

	// Step 5: Filter exchanges by multiple criteria
	fmt.Println("Step 5: Advanced filtering - Exchanges by area code and cost...")
	fmt.Println("---------------------------------------------------------------")

	if len(areaCodesResponse.Data) > 0 {
		firstAreaCode := areaCodesResponse.Data[0]
		acNum, err := strconv.Atoi(firstAreaCode.ID)
		if err == nil {
			filteredExchanges, err := client.Numbers().Purchaseable().
				ListAvailableExchanges(ctx).
				AreaCode(acNum).
				MaxSetupCost(2.0).
				Limit(10).
				Send()
			if err != nil {
				log.Printf("Error with filtered exchanges: %v", err)
			} else {
				fmt.Printf("Exchanges in area code %d with setup cost ≤ $2.00: %d\n",
					acNum, len(filteredExchanges.Data))

				maxToShow := 5
				if len(filteredExchanges.Data) < maxToShow {
					maxToShow = len(filteredExchanges.Data)
				}

				for i := 0; i < maxToShow; i++ {
					fmt.Printf("  • %s\n", filteredExchanges.Data[i].ID)
				}

				if len(filteredExchanges.Data) > 5 {
					fmt.Printf("  ... and %d more\n", len(filteredExchanges.Data)-5)
				}
			}
		}
	}
	fmt.Println()

	// Step 6: Search for actual phone numbers
	fmt.Println("Step 6: Searching for purchasable phone numbers...")
	fmt.Println("---------------------------------------------------")

	if len(areaCodesResponse.Data) > 0 {
		firstAreaCode := areaCodesResponse.Data[0]
		ac := firstAreaCode.ID
		fmt.Printf("\nSearching for numbers starting with %s...\n", ac)

		phoneNumbers, err := client.Numbers().Purchaseable().
			SearchAvailable(ctx).
			StartsWith(ac).
			Limit(5).
			Send()
		if err != nil {
			log.Printf("Error searching phone numbers: %v", err)
		} else {
			fmt.Printf("Found %d available phone numbers:\n", len(phoneNumbers.Data))
			for _, number := range phoneNumbers.Data {
				state := "N/A"
				if number.Attributes.State != nil {
					state = *number.Attributes.State
				}
				rateCenter := "N/A"
				if number.Attributes.RateCenter != nil {
					rateCenter = *number.Attributes.RateCenter
				}
				numberType := "N/A"
				if number.Attributes.NumberType != nil {
					numberType = *number.Attributes.NumberType
				}

				fmt.Printf("  • %s - $%.2f/month (Setup: $%.2f)\n",
					number.Attributes.Value,
					number.Attributes.MonthlyCost,
					number.Attributes.SetupCost)
				fmt.Printf("    State: %s, Rate Center: %s, Type: %s\n",
					state, rateCenter, numberType)
			}
		}
	}
	fmt.Println()

	// Step 7: Summary statistics
	fmt.Println("Step 7: Summary Statistics")
	fmt.Println("--------------------------")

	allAreaCodes, err := client.Numbers().Purchaseable().
		ListAvailableAreaCodes(ctx).
		Limit(200).
		Send()
	if err != nil {
		log.Printf("Error getting all area codes: %v", err)
	} else {
		fmt.Printf("Total area codes available: %d\n", len(allAreaCodes.Data))
	}

	allTiers, err := client.Numbers().Purchaseable().
		ListAvailableTiers(ctx).
		Limit(200).
		Send()
	if err != nil {
		log.Printf("Error getting all tiers: %v", err)
	} else {
		fmt.Printf("Total pricing tiers: %d\n", len(allTiers.Data))
	}

	if allAreaCodes != nil && len(allAreaCodes.Data) > 0 {
		firstAC := allAreaCodes.Data[0]
		acNum, err := strconv.Atoi(firstAC.ID)
		if err == nil {
			exchanges, err := client.Numbers().Purchaseable().
				ListAvailableExchanges(ctx).
				AreaCode(acNum).
				Limit(200).
				Send()
			if err != nil {
				log.Printf("Error getting exchanges: %v", err)
			} else {
				fmt.Printf("Sample exchanges in area code %d: %d\n", acNum, len(exchanges.Data))
			}
		}
	}
	fmt.Println()

	// Step 8: Display pagination information
	fmt.Println("Step 8: Pagination Links")
	fmt.Println("-----------------------")
	if allAreaCodes != nil {
		fmt.Println("Area Codes:")
		fmt.Printf("  Self: %s\n", allAreaCodes.Links.Self)
		if allAreaCodes.Links.Next != nil {
			fmt.Printf("  Next: %s\n", *allAreaCodes.Links.Next)
		} else {
			fmt.Println("  Next: None (all results retrieved)")
		}
	}

	if allTiers != nil {
		fmt.Println("\nTiers:")
		fmt.Printf("  Self: %s\n", allTiers.Links.Self)
		if allTiers.Links.Next != nil {
			fmt.Printf("  Next: %s\n", *allTiers.Links.Next)
		} else {
			fmt.Println("  Next: None (all results retrieved)")
		}
	}
	fmt.Println()

	// Step 9: Demonstrate builder pattern flexibility
	fmt.Println("Step 9: Builder Pattern Flexibility")
	fmt.Println("-----------------------------------")
	fmt.Println("Creating requests with different configurations:")
	fmt.Println()

	fmt.Println("→ Minimal configuration (defaults):")
	_ = client.Numbers().Purchaseable().ListAvailableAreaCodes(ctx)
	fmt.Println("  ✓ Created with default limit=10, offset=nil")
	fmt.Println()

	fmt.Println("→ Custom limit only:")
	_ = client.Numbers().Purchaseable().
		ListAvailableExchanges(ctx).
		Limit(25)
	fmt.Println("  ✓ Created with limit=25")
	fmt.Println()

	fmt.Println("→ Full configuration:")
	_ = client.Numbers().Purchaseable().
		ListAvailableExchanges(ctx).
		AreaCode(415).
		MaxSetupCost(5.0).
		Limit(50).
		Offset(10)
	fmt.Println("  ✓ Created with all parameters set")
	fmt.Println()

	fmt.Println("→ Phone number search:")
	_ = client.Numbers().Purchaseable().
		SearchAvailable(ctx).
		StartsWith("206").
		State("wa").
		NumberType("longcode").
		OrderBy("monthly_cost", "asc").
		Limit(10)
	fmt.Println("  ✓ Created with multiple search criteria")
	fmt.Println()

	// Final summary
	fmt.Println("✓ Comprehensive search demonstration completed successfully!")
	fmt.Println("\nKey Features Demonstrated:")
	fmt.Println("  • List available tiers")
	fmt.Println("  • List available area codes")
	fmt.Println("  • List available exchanges")
	fmt.Println("  • Search for purchasable phone numbers")
	fmt.Println("  • Filter by max setup cost")
	fmt.Println("  • Filter by area code, state, rate center")
	fmt.Println("  • Filter by number type (toll-free/longcode)")
	fmt.Println("  • Pattern matching (starts_with, contains, ends_with)")
	fmt.Println("  • Ordering results")
	fmt.Println("  • Pagination support")
	fmt.Println("  • Builder pattern with method chaining")
	fmt.Println("  • Cross-endpoint data exploration")
	fmt.Println("  • Context support for cancellation and timeouts")
}
