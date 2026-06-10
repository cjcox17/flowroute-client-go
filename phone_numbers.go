package flowroute

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// PhoneNumbersClient provides methods for managing account phone numbers
type PhoneNumbersClient struct {
	client *Client
}

// List returns a builder for listing account phone numbers
func (p *PhoneNumbersClient) List(ctx context.Context) *ListPhoneNumbersBuilder {
	return &ListPhoneNumbersBuilder{
		client: p.client,
		ctx:    ctx,
		limit:  10, // default limit
	}
}

// Purchase returns a builder for purchasing a phone number
func (p *PhoneNumbersClient) Purchase(ctx context.Context, numberID string) *PurchasePhoneNumberBuilder {
	return &PurchasePhoneNumberBuilder{
		client:   p.client,
		ctx:      ctx,
		numberID: numberID,
	}
}

// RemovePhoneNumber removes (deletes) a phone number from your account
//
// Example:
//
//	err := client.Numbers().PhoneNumbers().RemovePhoneNumber(ctx, "12065551234")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (p *PhoneNumbersClient) RemovePhoneNumber(ctx context.Context, numberID string) error {
	path := fmt.Sprintf("/v2/numbers/%s", numberID)
	resp, err := p.client.doRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return parseErrorResponse(resp)
	}

	return nil
}

// ListPhoneNumbersResponse represents the response from listing account phone numbers
type ListPhoneNumbersResponse struct {
	Data  []AccountPhoneNumber       `json:"data"`
	Links PhoneNumberPaginationLinks `json:"links"`
}

// AccountPhoneNumber represents a phone number in your account
type AccountPhoneNumber struct {
	Attributes PhoneNumberAttributes `json:"attributes"`
	ID         string                `json:"id"`
	Links      PhoneNumberSelfLink   `json:"links"`
	Type       string                `json:"type"`
}

// PhoneNumberAttributes contains the attributes of a phone number
type PhoneNumberAttributes struct {
	Alias              *string `json:"alias,omitempty"`
	CNAMLookupsEnabled *bool   `json:"cnam_lookups_enabled,omitempty"`
	MessagingEnabled   *bool   `json:"messaging_enabled,omitempty"`
	NumberType         string  `json:"number_type"`
	RateCenter         string  `json:"rate_center"`
	State              string  `json:"state"`
	Status             *string `json:"status,omitempty"`
	ISOCountry         *string `json:"iso_country,omitempty"`
	Value              string  `json:"value"`
}

// PhoneNumberSelfLink contains self-referential links
type PhoneNumberSelfLink struct {
	Self string `json:"self"`
}

// PhoneNumberPaginationLinks contains pagination links
type PhoneNumberPaginationLinks struct {
	Next *string `json:"next,omitempty"`
	Self string  `json:"self"`
}

// ListPhoneNumbersBuilder builds a request to list account phone numbers
type ListPhoneNumbersBuilder struct {
	client     *Client
	ctx        context.Context
	startsWith *string
	contains   *string
	endsWith   *string
	alias      *string
	limit      int
	offset     *int
}

// StartsWith filters phone numbers starting with the specified value
//
// To specify a country, use + and country code (e.g., "+44"),
// otherwise a default country code of '1' will be used.
func (b *ListPhoneNumbersBuilder) StartsWith(value string) *ListPhoneNumbersBuilder {
	b.startsWith = &value
	return b
}

// Contains filters phone numbers containing the specified value
func (b *ListPhoneNumbersBuilder) Contains(value string) *ListPhoneNumbersBuilder {
	b.contains = &value
	return b
}

// EndsWith filters phone numbers ending with the specified value
func (b *ListPhoneNumbersBuilder) EndsWith(value string) *ListPhoneNumbersBuilder {
	b.endsWith = &value
	return b
}

// Alias filters by phone number alias (exact match only)
func (b *ListPhoneNumbersBuilder) Alias(value string) *ListPhoneNumbersBuilder {
	b.alias = &value
	return b
}

// Limit sets the maximum number of items to retrieve (default: 10, max: 200)
func (b *ListPhoneNumbersBuilder) Limit(limit int) *ListPhoneNumbersBuilder {
	b.limit = limit
	return b
}

// Offset sets the offset for pagination
func (b *ListPhoneNumbersBuilder) Offset(offset int) *ListPhoneNumbersBuilder {
	b.offset = &offset
	return b
}

// Send executes the request and returns the response
func (b *ListPhoneNumbersBuilder) Send() (*ListPhoneNumbersResponse, error) {
	query := url.Values{}

	if b.startsWith != nil {
		query.Set("starts_with", *b.startsWith)
	}
	if b.contains != nil {
		query.Set("contains", *b.contains)
	}
	if b.endsWith != nil {
		query.Set("ends_with", *b.endsWith)
	}
	if b.alias != nil {
		query.Set("alias", *b.alias)
	}
	query.Set("limit", strconv.Itoa(b.limit))
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	resp, err := b.client.doRequest(b.ctx, http.MethodGet, "/v2/numbers", query, nil)
	if err != nil {
		return nil, err
	}

	var result ListPhoneNumbersResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// PurchasePhoneNumberResponse represents the response from purchasing a phone number
type PurchasePhoneNumberResponse struct {
	Data     PurchasedPhoneNumber `json:"data"`
	Included []IncludedObject     `json:"included,omitempty"`
	Links    PhoneNumberSelfLink  `json:"links"`
}

// PurchasedPhoneNumber represents a newly purchased phone number
type PurchasedPhoneNumber struct {
	Attributes    PurchasedPhoneNumberAttributes `json:"attributes"`
	ID            string                         `json:"id"`
	Links         PhoneNumberSelfLink            `json:"links"`
	Relationships PhoneNumberRelationships       `json:"relationships"`
	Type          string                         `json:"type"`
}

// PurchasedPhoneNumberAttributes contains attributes of a purchased phone number
type PurchasedPhoneNumberAttributes struct {
	Alias              *string  `json:"alias,omitempty"`
	CNAMLookupsEnabled *bool    `json:"cnam_lookups_enabled,omitempty"`
	InboundRate        *float64 `json:"inbound_rate,omitempty"`
	ISOCountry         *string  `json:"iso_country,omitempty"`
	MessagingEnabled   *bool    `json:"messaging_enabled,omitempty"`
	MonthlyCost        *float64 `json:"monthly_cost,omitempty"`
	NumberType         string   `json:"number_type"`
	RateCenter         string   `json:"rate_center"`
	RateType           *string  `json:"rate_type,omitempty"`
	SetupCost          *float64 `json:"setup_cost,omitempty"`
	State              string   `json:"state"`
	Status             *string  `json:"status,omitempty"`
	Tier               *string  `json:"tier,omitempty"`
	Value              string   `json:"value"`
}

// PhoneNumberRelationships contains relationships to other resources
type PhoneNumberRelationships struct {
	CNAMPreset       *RelationshipData `json:"cnam_preset,omitempty"`
	E911Address      *RelationshipData `json:"e911_address,omitempty"`
	FailoverRoute    *RelationshipData `json:"failover_route,omitempty"`
	MessageProvision *RelationshipData `json:"message_provision,omitempty"`
	PrimaryRoute     *RelationshipData `json:"primary_route,omitempty"`
}

// RelationshipData wraps related object data
type RelationshipData struct {
	Data *RelatedObject `json:"data,omitempty"`
}

// RelatedObject represents a related resource
type RelatedObject struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// IncludedObject represents an included related object in the response
type IncludedObject struct {
	Type       string              `json:"type"`
	ID         string              `json:"id"`
	Attributes json.RawMessage     `json:"attributes"`
	Links      PhoneNumberSelfLink `json:"links"`
}

// RouteAttributes contains attributes of a route
type RouteAttributes struct {
	Alias          *string `json:"alias,omitempty"`
	EdgeStrategyID *int    `json:"edge_strategy_id,omitempty"`
	RouteType      string  `json:"route_type"`
	Value          *string `json:"value,omitempty"`
}

// MessageProvisioningAttributes contains message provisioning status
type MessageProvisioningAttributes struct {
	Status string `json:"status"`
	TN     string `json:"tn"`
}

// PurchasePhoneNumberBuilder builds a request to purchase a phone number
type PurchasePhoneNumberBuilder struct {
	client               *Client
	ctx                  context.Context
	numberID             string
	cnamLookupsEnabled   *bool
	messagingEnabled     *bool
	messagingEmail       *bool
	messagingCallbackURL *string
}

// CNAMLookupsEnabled enables or disables CNAM (Caller-ID Name) Lookup
func (b *PurchasePhoneNumberBuilder) CNAMLookupsEnabled(enabled bool) *PurchasePhoneNumberBuilder {
	b.cnamLookupsEnabled = &enabled
	return b
}

// MessagingEnabled enables or disables messaging for the phone number
func (b *PurchasePhoneNumberBuilder) MessagingEnabled(enabled bool) *PurchasePhoneNumberBuilder {
	b.messagingEnabled = &enabled
	return b
}

// MessagingEmail enables or disables email notification for messaging provisioning
func (b *PurchasePhoneNumberBuilder) MessagingEmail(enabled bool) *PurchasePhoneNumberBuilder {
	b.messagingEmail = &enabled
	return b
}

// MessagingCallbackURL sets a callback URL for messaging provisioning status
func (b *PurchasePhoneNumberBuilder) MessagingCallbackURL(url string) *PurchasePhoneNumberBuilder {
	b.messagingCallbackURL = &url
	return b
}

// Send executes the purchase request
func (b *PurchasePhoneNumberBuilder) Send() (*PurchasePhoneNumberResponse, error) {
	requestBody := map[string]interface{}{}

	if b.cnamLookupsEnabled != nil {
		requestBody["cnam_lookups_enabled"] = *b.cnamLookupsEnabled
	}
	if b.messagingEnabled != nil {
		requestBody["messaging_enabled"] = *b.messagingEnabled
	}
	if b.messagingEmail != nil {
		requestBody["messaging_email"] = *b.messagingEmail
	}
	if b.messagingCallbackURL != nil {
		requestBody["messaging_callback_url"] = *b.messagingCallbackURL
	}

	path := fmt.Sprintf("/v2/numbers/%s", b.numberID)
	resp, err := b.client.doRequest(b.ctx, http.MethodPost, path, nil, requestBody)
	if err != nil {
		return nil, err
	}

	var result PurchasePhoneNumberResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
