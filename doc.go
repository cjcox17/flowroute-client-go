/*
Package flowroute provides a Go client library for the Flowroute API.

The Flowroute API allows you to manage phone numbers, routes, and messaging
services programmatically. This client provides a type-safe, idiomatic Go
interface to the API with support for context-based cancellation and timeouts.

# Installation

	go get github.com/cjcox17/flowroute-client-go

# Quick Start

Create a client and start making API calls:

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

		// Search for available phone numbers
		resp, err := client.Numbers().Purchaseable().
			SearchAvailable(ctx).
			StartsWith("206").
			State("WA").
			Limit(10).
			Send()
		if err != nil {
			log.Fatal(err)
		}

		for _, number := range resp.Data {
			fmt.Printf("%s - $%.2f/month\n",
				number.Attributes.Value,
				number.Attributes.MonthlyCost)
		}
	}

# Client Configuration

Create a client with custom options:

	client := flowroute.NewClient("api-key", "api-secret",
		flowroute.WithTimeout(60 * time.Second),
	)

# Context Support

All API calls accept a context for cancellation and timeouts:

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Numbers().Purchaseable().
		ListAvailableAreaCodes(ctx).
		Send()

# Error Handling

API errors can be checked using type assertions:

	import "errors"

	err := client.Numbers().PhoneNumbers().
		RemovePhoneNumber(ctx, "12065551234")
	if err != nil {
		var apiErr *flowroute.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("API Error: Status %d\n", apiErr.StatusCode)
		}
	}

# Builder Pattern

All requests use a fluent builder pattern for easy composition:

	resp, err := client.Numbers().Purchaseable().
		SearchAvailable(ctx).
		StartsWith("206").
		State("WA").
		NumberType("longcode").
		OrderBy("monthly_cost", "asc").
		Limit(10).
		Send()

# Features

  - Type-safe API interactions with comprehensive error handling
  - Context support for cancellation and timeouts
  - Builder pattern for complex request construction
  - Zero external dependencies (uses only Go standard library)
  - Secure Basic Authentication
  - Extensive documentation and examples

# API Coverage

This library provides access to the following Flowroute API endpoints:

Numbers API:
  - List available area codes
  - List available exchanges
  - List available pricing tiers
  - Search for purchasable phone numbers
  - List account phone numbers
  - Purchase phone numbers
  - Remove phone numbers from account

Messages API:
  - Send SMS and MMS
  - Look up a message detail record (MDR)
  - Look up a set of messages
  - Set account-level callback URLs (SMS, MMS, DLR)
  - Look up account-level callback details
  - Remove account-level callbacks
  - Parse inbound message and DLR webhook payloads

For more information, visit:
  - API Documentation: https://developer.flowroute.com/
  - Repository: https://github.com/cjcox17/flowroute-client-go
*/
package flowroute
