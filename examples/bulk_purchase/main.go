package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	flowroute "github.com/cjcox17/flowroute-client-go"
)

func main() {
	// Load .env file from current directory if present
	_ = loadEnv(".env")

	var (
		csvPath    = flag.String("csv", "", "Path to CSV file containing phone numbers to exclude (plain digits, no header)")
		startsWith = flag.String("starts-with", "", "Search prefix for available phone numbers")
		quantity   = flag.Int("quantity", 0, "Number of phone numbers to purchase")
		purchase   = flag.Bool("purchase", false, "Actually purchase numbers (default is dry-run)")
		maxPages   = flag.Int("max-pages", 50, "Maximum search result pages to fetch (safety limit)")
		pageLimit  = flag.Int("page-limit", 50, "Limit per search request")
	)
	flag.Parse()

	if *csvPath == "" {
		log.Fatal("-csv is required")
	}
	if *startsWith == "" {
		log.Fatal("-starts-with is required")
	}
	if *quantity <= 0 {
		log.Fatal("-quantity must be greater than 0")
	}

	// Load excluded numbers from CSV
	excluded, err := loadExcludedNumbers(*csvPath)
	if err != nil {
		log.Fatalf("Failed to load CSV: %v", err)
	}
	fmt.Printf("Loaded %d excluded numbers from %s\n", len(excluded), *csvPath)

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

	fmt.Printf("\nFlowroute Bulk Purchase Example\n")
	fmt.Println("================================")
	fmt.Printf("Search prefix: %s\n", *startsWith)
	fmt.Printf("Target quantity: %d\n", *quantity)
	if *purchase {
		fmt.Println("Mode: LIVE PURCHASE")
	} else {
		fmt.Println("Mode: DRY-RUN (use -purchase to actually buy)")
	}
	fmt.Println()

	// Search for available numbers, filtering out excluded ones
	var toPurchase []string
	offset := 0
	pages := 0

	for len(toPurchase) < *quantity && pages < *maxPages {
		pages++

		searchResp, err := client.Numbers().Purchaseable().
			SearchAvailable(ctx).
			StartsWith(*startsWith).
			Limit(*pageLimit).
			Offset(offset).
			Send()
		if err != nil {
			log.Fatalf("Search failed (page %d, offset %d): %v", pages, offset, err)
		}

		if len(searchResp.Data) == 0 {
			fmt.Println("No more results available.")
			break
		}

		for _, number := range searchResp.Data {
			value := normalizeNumber(number.Attributes.Value)
			if excluded[value] {
				continue
			}
			toPurchase = append(toPurchase, number.ID)
			if len(toPurchase) >= *quantity {
				break
			}
		}

		// If fewer results than the page limit were returned, we've reached the end
		if len(searchResp.Data) < *pageLimit {
			fmt.Println("Reached end of search results.")
			break
		}

		offset += *pageLimit
	}

	if len(toPurchase) == 0 {
		fmt.Println("No eligible numbers found to purchase.")
		return
	}

	fmt.Printf("\nFound %d eligible number(s) to purchase:\n", len(toPurchase))
	for i, id := range toPurchase {
		fmt.Printf("  %d. %s\n", i+1, id)
	}

	if !*purchase {
		fmt.Println("\nDry-run complete. No numbers were purchased.")
		fmt.Println("Re-run with -purchase to execute the purchase.")
		return
	}

	fmt.Println("\nProceeding with purchase...")

	var succeeded []string
	var failed []string

	for _, numberID := range toPurchase {
		purchaseCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		resp, err := client.Numbers().PhoneNumbers().
			Purchase(purchaseCtx, numberID).
			Send()
		cancel()

		if err != nil {
			log.Printf("  FAILED to purchase %s: %v", numberID, err)
			failed = append(failed, numberID)
			continue
		}

		fmt.Printf("  SUCCESS: Purchased %s\n", resp.Data.Attributes.Value)
		succeeded = append(succeeded, resp.Data.Attributes.Value)
	}

	// Write purchased numbers to purchased.csv
	if len(succeeded) > 0 {
		if err := writePurchasedCSV("purchased.csv", succeeded); err != nil {
			log.Printf("Warning: failed to write purchased.csv: %v", err)
		} else {
			fmt.Printf("Wrote %d purchased number(s) to purchased.csv\n", len(succeeded))
		}
	}

	fmt.Println("\nPurchase Summary")
	fmt.Println("----------------")
	fmt.Printf("Succeeded: %d\n", len(succeeded))
	fmt.Printf("Failed:    %d\n", len(failed))

	if len(failed) > 0 {
		fmt.Println("\nFailed numbers:")
		for _, id := range failed {
			fmt.Printf("  - %s\n", id)
		}
	}
}

// writePurchasedCSV writes a list of phone numbers to a CSV file (no header).
func writePurchasedCSV(path string, numbers []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	for _, number := range numbers {
		if err := writer.Write([]string{number}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

// loadExcludedNumbers reads a CSV file (no header, one number per line) and
// returns a set of normalized phone numbers.
func loadExcludedNumbers(path string) (map[string]bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(bufio.NewReader(f))
	// We expect a single column per row, but allow for extra columns just in case
	reader.FieldsPerRecord = -1

	excluded := make(map[string]bool)
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		if len(record) == 0 {
			continue
		}
		val := normalizeNumber(record[0])
		if val != "" {
			excluded[val] = true
		}
	}

	return excluded, nil
}

// normalizeNumber strips common formatting characters and whitespace from a
// phone number, returning a clean digit string.
func normalizeNumber(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, strings.TrimSpace(s))
}

// loadEnv reads a simple KEY=VALUE .env file and sets the variables via
// os.Setenv. Lines starting with # or that are empty are ignored. If the file
// does not exist, no error is returned.
func loadEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		// Strip optional surrounding quotes
		value = strings.Trim(value, "\"'")
		if key != "" {
			_ = os.Setenv(key, value)
		}
	}
	return scanner.Err()
}
