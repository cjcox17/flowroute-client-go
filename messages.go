package flowroute

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// MessagesClient provides methods for the Flowroute Messaging API

type MessagesClient struct {
	client *Client
}

// Send returns a builder for sending an SMS or MMS
func (m *MessagesClient) Send(ctx context.Context) *SendMessageBuilder {
	return &SendMessageBuilder{
		client: m.client,
		ctx:    ctx,
	}
}

// Get returns a builder for looking up a message detail record
func (m *MessagesClient) Get(ctx context.Context, recordID string) *GetMessageBuilder {
	return &GetMessageBuilder{
		client:   m.client,
		ctx:      ctx,
		recordID: recordID,
	}
}

// List returns a builder for retrieving a set of messages
func (m *MessagesClient) List(ctx context.Context) *ListMessagesBuilder {
	return &ListMessagesBuilder{
		client: m.client,
		ctx:    ctx,
		limit:  10,
	}
}

// SetCallback returns a builder for setting an account-level callback URL
func (m *MessagesClient) SetCallback(ctx context.Context, callbackType CallbackType) *SetCallbackBuilder {
	return &SetCallbackBuilder{
		client:       m.client,
		ctx:          ctx,
		callbackType: string(callbackType),
	}
}

// GetCallback returns a builder for looking up an account-level callback
func (m *MessagesClient) GetCallback(ctx context.Context, callbackType CallbackType) *GetCallbackBuilder {
	return &GetCallbackBuilder{
		client:       m.client,
		ctx:          ctx,
		callbackType: string(callbackType),
	}
}

// RemoveCallback removes an account-level callback URL
func (m *MessagesClient) RemoveCallback(ctx context.Context, callbackType CallbackType) error {
	path := fmt.Sprintf("/v2.2/messages/%s", url.PathEscape(string(callbackType)))
	resp, err := m.client.doRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return parseErrorResponse(resp)
	}
	return nil
}

// CallbackType represents the supported callback types

type CallbackType string

// Supported callback types
const (
	SMSCallback     CallbackType = "sms_callback"
	MMSCallback     CallbackType = "mms_callback"
	SMSDLRCallback  CallbackType = "sms_dlr_callback"
	MMSDLRCallback  CallbackType = "mms_dlr_callback"
)

// SendMessageBuilder builds a request to send an SMS or MMS

type SendMessageBuilder struct {
	client       *Client
	ctx          context.Context
	to           string
	from         string
	body         string
	isMMS        *bool
	mediaURLs    []string
	dlrCallback  *string
}

// To sets the destination phone number (E.164 format)
func (b *SendMessageBuilder) To(number string) *SendMessageBuilder {
	b.to = number
	return b
}

// From sets the source Flowroute phone number (E.164 format)
func (b *SendMessageBuilder) From(number string) *SendMessageBuilder {
	b.from = number
	return b
}

// Body sets the message content
func (b *SendMessageBuilder) Body(text string) *SendMessageBuilder {
	b.body = text
	return b
}

// IsMMS marks the message as an MMS (useful for text-only MMS)
func (b *SendMessageBuilder) IsMMS(isMMS bool) *SendMessageBuilder {
	b.isMMS = &isMMS
	return b
}

// MediaURLs attaches media files to an MMS
func (b *SendMessageBuilder) MediaURLs(urls ...string) *SendMessageBuilder {
	b.mediaURLs = append(b.mediaURLs, urls...)
	return b
}

// DLRCallback sets a per-message delivery receipt callback URL
func (b *SendMessageBuilder) DLRCallback(url string) *SendMessageBuilder {
	b.dlrCallback = &url
	return b
}

// Send executes the send request and returns the response
func (b *SendMessageBuilder) Send() (*SendMessageResponse, error) {
	if b.to == "" || b.from == "" {
		return nil, fmt.Errorf("to and from are required")
	}
	if b.body == "" && len(b.mediaURLs) == 0 {
		return nil, fmt.Errorf("body or media_urls are required")
	}

	attrs := map[string]interface{}{
		"to":   b.to,
		"from": b.from,
	}
	if b.body != "" {
		attrs["body"] = b.body
	}
	if b.isMMS != nil {
		attrs["is_mms"] = *b.isMMS
	}
	if len(b.mediaURLs) > 0 {
		attrs["media_urls"] = b.mediaURLs
	}
	if b.dlrCallback != nil {
		attrs["dlr_callback"] = *b.dlrCallback
	}

	body := map[string]interface{}{
		"data": map[string]interface{}{
			"type":       "message",
			"attributes": attrs,
		},
	}

	resp, err := b.client.doRequest(b.ctx, http.MethodPost, "/v2.2/messages", nil, body)
	if err != nil {
		return nil, err
	}

	var result SendMessageResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SendMessageResponse represents the response from sending a message

type SendMessageResponse struct {
	Data  MessageData `json:"data"`
	Links MessageLinks `json:"links,omitempty"`
}

// MessageData contains the core message data

type MessageData struct {
	ID           string            `json:"id"`
	Attributes   MessageAttributes `json:"attributes"`
	Links        MessageLinks      `json:"links,omitempty"`
	Type         string            `json:"type"`
	Relationships *MessageRelationships `json:"relationships,omitempty"`
	Included     []MediaObject         `json:"included,omitempty"`
}

// MessageAttributes contains the attributes of a message

type MessageAttributes struct {
	AmountDisplay       string             `json:"amount_display,omitempty"`
	AmountNanodollars   int64              `json:"amount_nanodollars,omitempty"`
	Body                string             `json:"body,omitempty"`
	DeliveryReceipts    []DeliveryReceipt  `json:"delivery_receipts,omitempty"`
	Direction           string             `json:"direction,omitempty"`
	From                string             `json:"from,omitempty"`
	IsMMS               bool               `json:"is_mms,omitempty"`
	MessageCallbackURL  string             `json:"message_callback_url,omitempty"`
	MessageEncoding     int                `json:"message_encoding,omitempty"`
	MessageType         string             `json:"message_type,omitempty"`
	PriceDetails        *PriceDetails      `json:"price_details,omitempty"`
	Status              string             `json:"status,omitempty"`
	Timestamp           string             `json:"timestamp,omitempty"`
	To                  string             `json:"to,omitempty"`
}

// DeliveryReceipt represents a delivery receipt for a message

type DeliveryReceipt struct {
	Level                   int     `json:"level"`
	Status                  string  `json:"status"`
	StatusCode              *int    `json:"status_code,omitempty"`
	StatusCodeDescription   string  `json:"status_code_description,omitempty"`
	Timestamp               string  `json:"timestamp"`
}

// PriceDetails contains the pricing breakdown of a message

type PriceDetails struct {
	BaseRate       string `json:"base_rate,omitempty"`
	ChargedCost    string `json:"charged_cost,omitempty"`
	SegmentCount   int    `json:"segment_count,omitempty"`
	SurchargeRate  string `json:"surcharge_rate,omitempty"`
}

// MessageLinks contains links related to a message

type MessageLinks struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
}

// MessageRelationships contains relationships for MMS media

type MessageRelationships struct {
	Media *MediaRelationship `json:"media,omitempty"`
}

// MediaRelationship contains the media data relationship

type MediaRelationship struct {
	Data []MediaRelationshipData `json:"data,omitempty"`
}

// MediaRelationshipData represents a single media relationship entry

type MediaRelationshipData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// MediaObject represents an included media attachment in an MMS

type MediaObject struct {
	Attributes MediaAttributes `json:"attributes"`
	ID         string          `json:"id"`
	Links      MediaLinks      `json:"links,omitempty"`
	Type       string          `json:"type"`
}

// MediaAttributes contains the attributes of a media attachment

type MediaAttributes struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	MIMEType string `json:"mime_type"`
	URL      string `json:"url"`
}

// MediaLinks contains links for a media object

type MediaLinks struct {
	Self string `json:"self,omitempty"`
}

// InboundMessage represents the payload posted by Flowroute to a callback URL
// for inbound SMS or MMS. Use this type to decode the JSON body in your HTTP handler.

type InboundMessage struct {
	Data MessageData `json:"data"`
}

// GetMessageBuilder builds a request to look up a single message detail record

type GetMessageBuilder struct {
	client   *Client
	ctx      context.Context
	recordID string
}

// Send executes the lookup request
func (b *GetMessageBuilder) Send() (*MessageData, error) {
	path := fmt.Sprintf("/v2.2/messages/%s", url.PathEscape(b.recordID))
	resp, err := b.client.doRequest(b.ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data MessageData `json:"data"`
	}
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// ListMessagesBuilder builds a request to retrieve a set of messages

type ListMessagesBuilder struct {
	client    *Client
	ctx       context.Context
	startDate *time.Time
	endDate   *time.Time
	limit     int
	offset    *int
}

// StartDate sets the beginning of the date range (UTC)
func (b *ListMessagesBuilder) StartDate(t time.Time) *ListMessagesBuilder {
	b.startDate = &t
	return b
}

// EndDate sets the end of the date range (UTC)
func (b *ListMessagesBuilder) EndDate(t time.Time) *ListMessagesBuilder {
	b.endDate = &t
	return b
}

// Limit sets the maximum number of results to return
func (b *ListMessagesBuilder) Limit(limit int) *ListMessagesBuilder {
	b.limit = limit
	return b
}

// Offset sets the pagination offset
func (b *ListMessagesBuilder) Offset(offset int) *ListMessagesBuilder {
	b.offset = &offset
	return b
}

// Send executes the request and returns the response
func (b *ListMessagesBuilder) Send() (*MessageListResponse, error) {
	query := url.Values{}
	query.Set("limit", strconv.Itoa(b.limit))

	if b.startDate != nil {
		query.Set("start_date", b.startDate.Format(time.RFC3339Nano))
	}
	if b.endDate != nil {
		query.Set("end_date", b.endDate.Format(time.RFC3339Nano))
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	resp, err := b.client.doRequest(b.ctx, http.MethodGet, "/v2.2/messages", query, nil)
	if err != nil {
		return nil, err
	}

	var result MessageListResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// MessageListResponse represents the response from listing messages

type MessageListResponse struct {
	Data  []MessageData   `json:"data"`
	Links PaginationLinks `json:"links"`
}

// SetCallbackBuilder builds a request to set an account-level callback URL

type SetCallbackBuilder struct {
	client       *Client
	ctx          context.Context
	callbackType string
	callbackURL  string
}

// URL sets the callback URL
func (b *SetCallbackBuilder) URL(url string) *SetCallbackBuilder {
	b.callbackURL = url
	return b
}

// Send executes the request
func (b *SetCallbackBuilder) Send() (*CallbackResponse, error) {
	if b.callbackURL == "" {
		return nil, fmt.Errorf("callback URL is required")
	}

	body := map[string]interface{}{
		"callback_url": b.callbackURL,
	}

	path := fmt.Sprintf("/v2.2/messages/%s", url.PathEscape(b.callbackType))
	resp, err := b.client.doRequest(b.ctx, http.MethodPut, path, nil, body)
	if err != nil {
		return nil, err
	}

	var result CallbackResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCallbackBuilder builds a request to look up an account-level callback

type GetCallbackBuilder struct {
	client       *Client
	ctx          context.Context
	callbackType string
}

// Send executes the request
func (b *GetCallbackBuilder) Send() (*CallbackResponse, error) {
	path := fmt.Sprintf("/v2.2/messages/%s", url.PathEscape(b.callbackType))
	resp, err := b.client.doRequest(b.ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result CallbackResponse
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CallbackResponse represents the response from callback management endpoints

type CallbackResponse struct {
	Data CallbackData `json:"data"`
}

// CallbackData contains the callback information

type CallbackData struct {
	Attributes CallbackAttributes `json:"attributes"`
	Type       string             `json:"type"`
}

// CallbackAttributes contains the callback attributes

type CallbackAttributes struct {
	CallbackURL string `json:"callback_url"`
	Product     string `json:"product,omitempty"`
}

// DeliveryReceiptPayload represents the payload posted by Flowroute to a
// DLR callback URL. Use this type to decode the JSON body in your HTTP handler.

type DeliveryReceiptPayload struct {
	Data DeliveryReceiptData `json:"data"`
}

// DeliveryReceiptData contains the core delivery receipt data

type DeliveryReceiptData struct {
	Attributes DeliveryReceipt `json:"attributes"`
	ID         string          `json:"id"`
	Links      MessageLinks    `json:"links,omitempty"`
	Type       string          `json:"type"`
}
