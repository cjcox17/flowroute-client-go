package flowroute

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// PurchaseableNumbersClient provides methods for searching purchaseable phone numbers
type PurchaseableNumbersClient struct {
	client *Client
}

// ListAvailableAreaCodes returns a builder for listing available area codes
func (p *PurchaseableNumbersClient) ListAvailableAreaCodes(ctx context.Context) *ListAreaCodesBuilder {
	return &ListAreaCodesBuilder{
		client: p.client,
		ctx:    ctx,
		limit:  10, // default
	}
}

// ListAvailableExchanges returns a builder for listing available exchanges
func (p *PurchaseableNumbersClient) ListAvailableExchanges(ctx context.Context) *ListExchangesBuilder {
	return &ListExchangesBuilder{
		client: p.client,
		ctx:    ctx,
		limit:  10, // default
	}
}

// ListAvailableTiers returns a builder for listing available tiers
func (p *PurchaseableNumbersClient) ListAvailableTiers(ctx context.Context) *ListTiersBuilder {
	return &ListTiersBuilder{
		client: p.client,
		ctx:    ctx,
		limit:  10, // default
	}
}

// SearchAvailable returns a builder for searching available phone numbers
func (p *PurchaseableNumbersClient) SearchAvailable(ctx context.Context) *SearchAvailableNumbersBuilder {
	return &SearchAvailableNumbersBuilder{
		client: p.client,
		ctx:    ctx,
		limit:  10, // default
	}
}

// PaginationLinks contains pagination links for list responses
type PaginationLinks struct {
	Next *string `json:"next,omitempty"`
	Self string  `json:"self"`
}

// AreaCodesResponse represents the response from listing area codes
type AreaCodesResponse struct {
	Data  []AreaCodeItem  `json:"data"`
	Links PaginationLinks `json:"links"`
}

// AreaCodeItem represents an area code
type AreaCodeItem struct {
	ID    string        `json:"id"`
	Links AreaCodeLinks `json:"links"`
	Type  string        `json:"type"`
}

// AreaCodeLinks contains links related to an area code
type AreaCodeLinks struct {
	Related string `json:"related"`
}

// ExchangesResponse represents the response from listing exchanges
type ExchangesResponse struct {
	Data  []ExchangeItem  `json:"data"`
	Links PaginationLinks `json:"links"`
}

// ExchangeItem represents an exchange
type ExchangeItem struct {
	ID    string        `json:"id"`
	Links ExchangeLinks `json:"links"`
	Type  string        `json:"type"`
}

// ExchangeLinks contains links related to an exchange
type ExchangeLinks struct {
	Related string `json:"related"`
}

// TiersResponse represents the response from listing tiers
type TiersResponse struct {
	Data  []TierItem      `json:"data"`
	Links PaginationLinks `json:"links"`
}

// TierItem represents a pricing tier
type TierItem struct {
	ID    string    `json:"id"`
	Links TierLinks `json:"links"`
	Type  string    `json:"type"`
}

// TierLinks contains links related to a tier
type TierLinks struct {
	Related string `json:"related"`
}

// SearchAvailableNumbersResponse represents the response from searching available numbers
type SearchAvailableNumbersResponse struct {
	Data  []AvailablePhoneNumber `json:"data"`
	Links PaginationLinks        `json:"links"`
}

// AvailablePhoneNumber represents an available phone number
type AvailablePhoneNumber struct {
	ID         string                    `json:"id"`
	Attributes AvailablePhoneNumberAttrs `json:"attributes"`
	Links      AvailablePhoneNumberLinks `json:"links"`
	Type       string                    `json:"type"`
}

// AvailablePhoneNumberAttrs contains attributes of an available phone number
type AvailablePhoneNumberAttrs struct {
	MonthlyCost float64 `json:"monthly_cost"`
	NumberType  *string `json:"number_type,omitempty"`
	RateCenter  *string `json:"rate_center,omitempty"`
	InboundRate float64 `json:"inbound_rate"`
	RateType    *string `json:"rate_type,omitempty"`
	Tier        *string `json:"tier,omitempty"`
	SetupCost   float64 `json:"setup_cost"`
	State       *string `json:"state,omitempty"`
	ISOCountry  *string `json:"iso_country,omitempty"`
	Value       string  `json:"value"`
}

// AvailablePhoneNumberLinks contains links for an available phone number
type AvailablePhoneNumberLinks struct {
	Related string `json:"related"`
}

// ListAreaCodesBuilder builds a request to list available area codes
type ListAreaCodesBuilder struct {
	client       *Client
	ctx          context.Context
	limit        int
	offset       *int
	maxSetupCost *float64
}

// Limit sets the maximum number of results to return
func (b *ListAreaCodesBuilder) Limit(limit int) *ListAreaCodesBuilder {
	b.limit = limit
	return b
}

// Offset sets the pagination offset
func (b *ListAreaCodesBuilder) Offset(offset int) *ListAreaCodesBuilder {
	b.offset = &offset
	return b
}

// MaxSetupCost filters by maximum setup cost
func (b *ListAreaCodesBuilder) MaxSetupCost(cost float64) *ListAreaCodesBuilder {
	b.maxSetupCost = &cost
	return b
}

// Send executes the request and returns the response
func (b *ListAreaCodesBuilder) Send() (*AreaCodesResponse, error) {
	query := url.Values{}
	query.Set("limit", strconv.Itoa(b.limit))

	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}
	if b.maxSetupCost != nil {
		query.Set("max_setup_cost", fmt.Sprintf("%.2f", *b.maxSetupCost))
	}

	resp, err := b.client.doRequest(b.ctx, http.MethodGet, "/v2/numbers/available/areacodes", query, nil)
	if err != nil {
		return nil, err
	}

	var result AreaCodesResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListExchangesBuilder builds a request to list available exchanges
type ListExchangesBuilder struct {
	client       *Client
	ctx          context.Context
	limit        int
	offset       *int
	maxSetupCost *float64
	areaCode     *int
}

// Limit sets the maximum number of results to return
func (b *ListExchangesBuilder) Limit(limit int) *ListExchangesBuilder {
	b.limit = limit
	return b
}

// Offset sets the pagination offset
func (b *ListExchangesBuilder) Offset(offset int) *ListExchangesBuilder {
	b.offset = &offset
	return b
}

// MaxSetupCost filters by maximum setup cost
func (b *ListExchangesBuilder) MaxSetupCost(cost float64) *ListExchangesBuilder {
	b.maxSetupCost = &cost
	return b
}

// AreaCode filters exchanges by area code
func (b *ListExchangesBuilder) AreaCode(areaCode int) *ListExchangesBuilder {
	b.areaCode = &areaCode
	return b
}

// Send executes the request and returns the response
func (b *ListExchangesBuilder) Send() (*ExchangesResponse, error) {
	query := url.Values{}
	query.Set("limit", strconv.Itoa(b.limit))

	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}
	if b.maxSetupCost != nil {
		query.Set("max_setup_cost", fmt.Sprintf("%.2f", *b.maxSetupCost))
	}
	if b.areaCode != nil {
		query.Set("areacode", strconv.Itoa(*b.areaCode))
	}

	resp, err := b.client.doRequest(b.ctx, http.MethodGet, "/v2/numbers/available/exchanges", query, nil)
	if err != nil {
		return nil, err
	}

	var result ExchangesResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListTiersBuilder builds a request to list available tiers
type ListTiersBuilder struct {
	client *Client
	ctx    context.Context
	limit  int
	offset *int
}

// Limit sets the maximum number of results to return
func (b *ListTiersBuilder) Limit(limit int) *ListTiersBuilder {
	b.limit = limit
	return b
}

// Offset sets the pagination offset
func (b *ListTiersBuilder) Offset(offset int) *ListTiersBuilder {
	b.offset = &offset
	return b
}

// Send executes the request and returns the response
func (b *ListTiersBuilder) Send() (*TiersResponse, error) {
	query := url.Values{}
	query.Set("limit", strconv.Itoa(b.limit))

	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	resp, err := b.client.doRequest(b.ctx, http.MethodGet, "/v2/numbers/available/tiers", query, nil)
	if err != nil {
		return nil, err
	}

	var result TiersResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SearchAvailableNumbersBuilder builds a request to search available phone numbers
type SearchAvailableNumbersBuilder struct {
	client     *Client
	ctx        context.Context
	limit      int
	offset     *int
	startsWith *string
	contains   *string
	endsWith   *string
	rateCenter *string
	state      *string
	numberType *string
	isoCountry *string
	orderBy    *string
}

// Limit sets the maximum number of results to return
func (b *SearchAvailableNumbersBuilder) Limit(limit int) *SearchAvailableNumbersBuilder {
	b.limit = limit
	return b
}

// Offset sets the pagination offset
func (b *SearchAvailableNumbersBuilder) Offset(offset int) *SearchAvailableNumbersBuilder {
	b.offset = &offset
	return b
}

// StartsWith filters numbers starting with the specified value
func (b *SearchAvailableNumbersBuilder) StartsWith(value string) *SearchAvailableNumbersBuilder {
	b.startsWith = &value
	return b
}

// Contains filters numbers containing the specified value
func (b *SearchAvailableNumbersBuilder) Contains(value string) *SearchAvailableNumbersBuilder {
	b.contains = &value
	return b
}

// EndsWith filters numbers ending with the specified value
func (b *SearchAvailableNumbersBuilder) EndsWith(value string) *SearchAvailableNumbersBuilder {
	b.endsWith = &value
	return b
}

// RateCenter filters by rate center
func (b *SearchAvailableNumbersBuilder) RateCenter(value string) *SearchAvailableNumbersBuilder {
	b.rateCenter = &value
	return b
}

// State filters by state (2-letter code)
func (b *SearchAvailableNumbersBuilder) State(value string) *SearchAvailableNumbersBuilder {
	b.state = &value
	return b
}

// NumberType filters by number type (e.g., "longcode", "toll-free")
func (b *SearchAvailableNumbersBuilder) NumberType(value string) *SearchAvailableNumbersBuilder {
	b.numberType = &value
	return b
}

// ISOCountry filters by ISO country code
func (b *SearchAvailableNumbersBuilder) ISOCountry(value string) *SearchAvailableNumbersBuilder {
	b.isoCountry = &value
	return b
}

// OrderBy sets the ordering of results
//
// Example:
//
//	builder.OrderBy("monthly_cost", "asc")
func (b *SearchAvailableNumbersBuilder) OrderBy(field, direction string) *SearchAvailableNumbersBuilder {
	orderBy := fmt.Sprintf("%s,%s", field, direction)
	b.orderBy = &orderBy
	return b
}

// Send executes the request and returns the response
func (b *SearchAvailableNumbersBuilder) Send() (*SearchAvailableNumbersResponse, error) {
	query := url.Values{}
	query.Set("limit", strconv.Itoa(b.limit))

	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}
	if b.startsWith != nil {
		query.Set("starts_with", *b.startsWith)
	}
	if b.contains != nil {
		query.Set("contains", *b.contains)
	}
	if b.endsWith != nil {
		query.Set("ends_with", *b.endsWith)
	}
	if b.rateCenter != nil {
		query.Set("rate_center", *b.rateCenter)
	}
	if b.state != nil {
		query.Set("state", *b.state)
	}
	if b.numberType != nil {
		query.Set("number_type", *b.numberType)
	}
	if b.isoCountry != nil {
		query.Set("iso_country", *b.isoCountry)
	}
	if b.orderBy != nil {
		query.Set("order_by", *b.orderBy)
	}

	resp, err := b.client.doRequest(b.ctx, http.MethodGet, "/v2/numbers/available", query, nil)
	if err != nil {
		return nil, err
	}

	var result SearchAvailableNumbersResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
