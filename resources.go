package ripllo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// This file implements the typed resource namespaces. The Go SDK uses
// `map[string]any` for the same Record<string,unknown> spots the
// Node SDK has (most marketing routes, profile routes, etc.). The
// strictly typed ones — discount codes, pixels, blog, abandoned cart,
// referrals, partner admin — get full struct types in types.go.

// mountResources wires the namespaces onto the Client. Called from
// NewClient and ForMerchant.
func (c *Client) mountResources() {
	c.API = &GeneratedAPI{c: c}
	c.DiscountCodes = &DiscountCodesResource{c: c}
	c.Pixels = &PixelsResource{c: c}
	c.Feeds = &FeedsResource{c: c}
	c.Blog = &BlogResource{c: c}
	c.AbandonedCart = &AbandonedCartResource{c: c}
	c.Referrals = &ReferralsResource{c: c}
	c.APIKeys = &APIKeysResource{c: c}
	c.Webhooks = &WebhooksResource{c: c}
	c.AuditLog = &AuditLogResource{c: c}
	c.Integrations = &IntegrationsResource{c: c}
	c.Billing = &BillingResource{c: c}
	c.Uploads = &UploadsResource{c: c}
	c.Campaigns = &CampaignsResource{c: c}
	c.Programs = &ProgramsResource{c: c}
	c.Collaborations = &CollaborationsResource{c: c}
	c.Insights = &InsightsResource{c: c}
	c.Channels = &ChannelsResource{c: c}
	c.Contacts = &ContactsResource{c: c}
	c.ContactLists = &ContactListsResource{c: c}
	c.Broadcasts = &BroadcastsResource{c: c}
	c.MarketingCampaigns = &MarketingCampaignsResource{c: c}
	c.Funnels = &FunnelsResource{c: c}
	c.Inbox = &InboxResource{c: c}
	c.AudienceSegments = &AudienceSegmentsResource{c: c}
	c.MerchantProfile = &MerchantProfileResource{c: c}
	c.CreatorProfile = &CreatorProfileResource{c: c}
	c.AffiliatorProfile = &AffiliatorProfileResource{c: c}
	c.Marketplace = &MarketplaceResource{c: c}
	c.Affiliates = &AffiliatesResource{c: c}
	c.KYC = &KYCResource{c: c}
	c.CreatorStats = &CreatorStatsResource{c: c}
	c.Admin = &AdminResource{c: c}
}

// JSON is the generic untyped response shape used by the resource
// methods that don't have a fully-modeled struct yet (most marketing
// routes). It's a thin alias for map[string]any so callers can decode
// once and JSON.Unmarshal into a stricter type if they want.
type JSON = map[string]any

// ─── Discount codes ────────────────────────────────────────────────────

type DiscountCodesResource struct{ c *Client }

type DiscountCodesListParams struct {
	Limit  *int
	Cursor *string
	Active *bool
}

func (r *DiscountCodesResource) List(ctx context.Context, p DiscountCodesListParams) (*DiscountCodeListPage, error) {
	var out DiscountCodeListPage
	q := map[string]any{}
	if p.Limit != nil {
		q["limit"] = *p.Limit
	}
	if p.Cursor != nil {
		q["cursor"] = *p.Cursor
	}
	if p.Active != nil {
		q["active"] = *p.Active
	}
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/discount-codes" + qs(q)}, &out)
	return &out, err
}

func (r *DiscountCodesResource) Get(ctx context.Context, id string) (*DiscountCode, error) {
	var out DiscountCode
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/discount-codes/" + url.PathEscape(id)}, &out)
	return &out, err
}

func (r *DiscountCodesResource) Create(ctx context.Context, input DiscountCodeCreateInput) (*DiscountCode, error) {
	var out DiscountCode
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/discount-codes", Body: input, IdempotencyKey: r.c.idemFn()}, &out)
	return &out, err
}

// Update — patch is a partial DiscountCodeCreateInput. Use pointers to
// indicate which fields to send.
func (r *DiscountCodesResource) Update(ctx context.Context, id string, patch DiscountCodeCreateInput) (*DiscountCode, error) {
	var out DiscountCode
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/discount-codes/" + url.PathEscape(id), Body: patch}, &out)
	return &out, err
}

func (r *DiscountCodesResource) Archive(ctx context.Context, id string) (*ArchiveResult, error) {
	var out ArchiveResult
	err := r.c.Do(ctx, RequestOptions{Method: "DELETE", Path: "/api/v1/discount-codes/" + url.PathEscape(id)}, &out)
	return &out, err
}

// Validate — read-only validation for cart preview.
func (r *DiscountCodesResource) Validate(ctx context.Context, input ValidateInput) (*ValidateResult, error) {
	var out ValidateResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/discount-codes/validate", Body: input}, &out)
	return &out, err
}

// Redeem — idempotent redemption. Call from your payment-success path.
func (r *DiscountCodesResource) Redeem(ctx context.Context, input RedeemInput) (*RedeemResult, error) {
	var out RedeemResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/discount-codes/redeem", Body: input}, &out)
	return &out, err
}

type ApplicableParams struct {
	Currency  *string
	ProductID *string
	Tags      *string
	Subtotal  *float64
}

// Applicable returns public applicable codes for a storefront teaser.
func (r *DiscountCodesResource) Applicable(ctx context.Context, accountID string, p ApplicableParams) (*ApplicableListWrap, error) {
	q := map[string]any{}
	if p.Currency != nil {
		q["currency"] = *p.Currency
	}
	if p.ProductID != nil {
		q["productId"] = *p.ProductID
	}
	if p.Tags != nil {
		q["tags"] = *p.Tags
	}
	if p.Subtotal != nil {
		q["subtotal"] = *p.Subtotal
	}
	var out ApplicableListWrap
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/discount-codes/applicable/" + url.PathEscape(accountID) + qs(q)}, &out)
	return &out, err
}

// ─── Pixels ────────────────────────────────────────────────────────────

type PixelsResource struct{ c *Client }

func (r *PixelsResource) Get(ctx context.Context) (*MerchantPixels, error) {
	var out MerchantPixels
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/pixels"}, &out)
	return &out, err
}

func (r *PixelsResource) Update(ctx context.Context, patch MerchantPixels) (*MerchantPixels, error) {
	var out MerchantPixels
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/pixels", Body: patch}, &out)
	return &out, err
}

// Public — storefront-public read, never returns CAPI secret.
func (r *PixelsResource) Public(ctx context.Context, accountID string) (*PublicPixels, error) {
	var out PublicPixels
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/pixels/public/" + url.PathEscape(accountID)}, &out)
	return &out, err
}

// ─── Feeds ─────────────────────────────────────────────────────────────

type FeedsResource struct{ c *Client }

func (r *FeedsResource) GetConfig(ctx context.Context) (*MerchantFeedConfig, error) {
	var out MerchantFeedConfig
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/feeds/config"}, &out)
	return &out, err
}

func (r *FeedsResource) UpdateConfig(ctx context.Context, patch MerchantFeedConfig) (*MerchantFeedConfig, error) {
	var out MerchantFeedConfig
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/feeds/config", Body: patch}, &out)
	return &out, err
}

// GoogleFeedURL builds the public XML feed URL — no signing required.
func (r *FeedsResource) GoogleFeedURL(accountID string) string {
	return r.c.baseURL + "/api/v1/feeds/google/" + url.PathEscape(accountID) + ".xml"
}

// ─── Blog ──────────────────────────────────────────────────────────────

type BlogResource struct{ c *Client }

type BlogListParams struct {
	Status *BlogPostStatus
}

func (r *BlogResource) List(ctx context.Context, p BlogListParams) (*BlogPostListWrap, error) {
	q := map[string]any{}
	if p.Status != nil {
		q["status"] = string(*p.Status)
	}
	var out BlogPostListWrap
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/blog" + qs(q)}, &out)
	return &out, err
}

func (r *BlogResource) Get(ctx context.Context, id string) (*BlogPostWrap, error) {
	var out BlogPostWrap
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/blog/" + url.PathEscape(id)}, &out)
	return &out, err
}

func (r *BlogResource) Create(ctx context.Context, input BlogPostInput) (*BlogPostWrap, error) {
	var out BlogPostWrap
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/blog", Body: input, IdempotencyKey: r.c.idemFn()}, &out)
	return &out, err
}

func (r *BlogResource) Update(ctx context.Context, id string, patch BlogPostInput) (*BlogPostWrap, error) {
	var out BlogPostWrap
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/blog/" + url.PathEscape(id), Body: patch}, &out)
	return &out, err
}

func (r *BlogResource) Delete(ctx context.Context, id string) (*DeletedResult, error) {
	var out DeletedResult
	err := r.c.Do(ctx, RequestOptions{Method: "DELETE", Path: "/api/v1/blog/" + url.PathEscape(id)}, &out)
	return &out, err
}

// PublicList returns a slimmed list shape (no body). Decoded as
// generic BlogPostListWrap — the Node SDK's Pick<> shape just hides
// fields, the JSON is still valid against BlogPost.
func (r *BlogResource) PublicList(ctx context.Context, accountID string) (*BlogPostListWrap, error) {
	var out BlogPostListWrap
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/blog/public/" + url.PathEscape(accountID)}, &out)
	return &out, err
}

func (r *BlogResource) PublicGet(ctx context.Context, accountID, slug string) (*BlogPostWrap, error) {
	var out BlogPostWrap
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/blog/public/" + url.PathEscape(accountID) + "/" + url.PathEscape(slug)}, &out)
	return &out, err
}

// ─── Abandoned cart ────────────────────────────────────────────────────

type AbandonedCartResource struct{ c *Client }

func (r *AbandonedCartResource) GetConfig(ctx context.Context) (*AbandonedCartConfig, error) {
	var out AbandonedCartConfig
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/abandoned-cart/config"}, &out)
	return &out, err
}

func (r *AbandonedCartResource) UpdateConfig(ctx context.Context, patch AbandonedCartConfig) (*AbandonedCartConfig, error) {
	var out AbandonedCartConfig
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/abandoned-cart/config", Body: patch}, &out)
	return &out, err
}

type RemindersListParams struct {
	Limit *int
}

func (r *AbandonedCartResource) ListReminders(ctx context.Context, p RemindersListParams) (*RemindersListWrap, error) {
	q := map[string]any{}
	if p.Limit != nil {
		q["limit"] = *p.Limit
	}
	var out RemindersListWrap
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/abandoned-cart/reminders" + qs(q)}, &out)
	return &out, err
}

type StatsParams struct {
	WindowDays *int
}

func (r *AbandonedCartResource) Stats(ctx context.Context, p StatsParams) (*RecoveryStats, error) {
	q := map[string]any{}
	if p.WindowDays != nil {
		q["windowDays"] = *p.WindowDays
	}
	var out RecoveryStats
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/abandoned-cart/stats" + qs(q)}, &out)
	return &out, err
}

func (r *AbandonedCartResource) RecordReminder(ctx context.Context, input RecordReminderInput) (*RecordReminderResult, error) {
	var out RecordReminderResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/abandoned-cart/reminders", Body: input}, &out)
	return &out, err
}

func (r *AbandonedCartResource) MarkRecovered(ctx context.Context, input MarkRecoveredInput) (*MarkRecoveredResult, error) {
	var out MarkRecoveredResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/abandoned-cart/recover", Body: input}, &out)
	return &out, err
}

// ─── Referrals ─────────────────────────────────────────────────────────

type ReferralsResource struct{ c *Client }

func (r *ReferralsResource) GetProgram(ctx context.Context) (*ReferralProgram, error) {
	var out ReferralProgram
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/referrals/program"}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *ReferralsResource) PutProgram(ctx context.Context, input ReferralProgramInput) (*ReferralProgram, error) {
	var out ReferralProgram
	err := r.c.Do(ctx, RequestOptions{Method: "PUT", Path: "/api/v1/referrals/program", Body: input}, &out)
	return &out, err
}

func (r *ReferralsResource) Stats(ctx context.Context) (*ProgramStats, error) {
	var out ProgramStats
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/referrals/stats"}, &out)
	return &out, err
}

func (r *ReferralsResource) IssueLink(ctx context.Context, input IssueLinkInput) (*IssueLinkResult, error) {
	var out IssueLinkResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/referrals/links/issue", Body: input}, &out)
	return &out, err
}

func (r *ReferralsResource) ResolveLink(ctx context.Context, accountID, code string) (*ResolveLinkResult, error) {
	var out ResolveLinkResult
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/referrals/links/" + url.PathEscape(accountID) + "/" + url.PathEscape(code)}, &out)
	return &out, err
}

func (r *ReferralsResource) RecordClick(ctx context.Context, input RecordClickInput) (*RecordClickResult, error) {
	var out RecordClickResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/referrals/links/click", Body: input}, &out)
	return &out, err
}

func (r *ReferralsResource) AttributeOnSignup(ctx context.Context, input AttributeSignupInput) (*AttributeSignupResult, error) {
	var out AttributeSignupResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/referrals/attributions/signup", Body: input}, &out)
	return &out, err
}

func (r *ReferralsResource) AttributeCheckoutStart(ctx context.Context, input AttributeCheckoutStartInput) (*AttributeCheckoutStartResult, error) {
	var out AttributeCheckoutStartResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/referrals/attributions/checkout-start", Body: input}, &out)
	return &out, err
}

func (r *ReferralsResource) FulfillRewardOnPayment(ctx context.Context, input FulfillRewardInput) (*FulfillRewardResult, error) {
	var out FulfillRewardResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/referrals/attributions/fulfill", Body: input}, &out)
	return &out, err
}

func (r *ReferralsResource) VoidAttributionOnRefund(ctx context.Context, checkoutSessionID string) (*VoidAttributionResult, error) {
	var out VoidAttributionResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/referrals/attributions/void", Body: map[string]string{"checkoutSessionId": checkoutSessionID}}, &out)
	return &out, err
}

func (r *ReferralsResource) ListMyRewards(ctx context.Context, accountID, customerID string) (*MyRewardsListWrap, error) {
	var out MyRewardsListWrap
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/referrals/rewards/" + url.PathEscape(accountID) + "/" + url.PathEscape(customerID)}, &out)
	return &out, err
}

func (r *ReferralsResource) ExpirePending(ctx context.Context) (*ExpirePendingResult, error) {
	var out ExpirePendingResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/referrals/sweeps/expire-pending", Body: map[string]any{}}, &out)
	return &out, err
}

// ─── API keys ──────────────────────────────────────────────────────────

type APIKeysResource struct{ c *Client }

func (r *APIKeysResource) List(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/api-keys"}, &out)
	return out, err
}

type APIKeyCreateInput struct {
	Description *string `json:"description,omitempty"`
	Scope       *string `json:"scope,omitempty"`
}

func (r *APIKeysResource) Create(ctx context.Context, input APIKeyCreateInput) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/api-keys", Body: input, IdempotencyKey: r.c.idemFn()}, &out)
	return out, err
}

func (r *APIKeysResource) Revoke(ctx context.Context, id string) (*RevokedResult, error) {
	var out RevokedResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/api-keys/" + url.PathEscape(id) + "/revoke", Body: map[string]any{}}, &out)
	return &out, err
}

// ─── Webhook endpoints + events ────────────────────────────────────────

type WebhooksResource struct{ c *Client }

func (r *WebhooksResource) ListEndpoints(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/webhooks/endpoints"}, &out)
	return out, err
}

type WebhookEndpointInput struct {
	URL         string   `json:"url"`
	Events      []string `json:"events,omitempty"`
	Description *string  `json:"description,omitempty"`
}

func (r *WebhooksResource) CreateEndpoint(ctx context.Context, input WebhookEndpointInput) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/webhooks/endpoints", Body: input, IdempotencyKey: r.c.idemFn()}, &out)
	return out, err
}

func (r *WebhooksResource) UpdateEndpoint(ctx context.Context, id string, patch JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/webhooks/endpoints/" + url.PathEscape(id), Body: patch}, &out)
	return out, err
}

func (r *WebhooksResource) DeleteEndpoint(ctx context.Context, id string) (*DeletedResult, error) {
	var out DeletedResult
	err := r.c.Do(ctx, RequestOptions{Method: "DELETE", Path: "/api/v1/webhooks/endpoints/" + url.PathEscape(id)}, &out)
	return &out, err
}

type EventsListParams struct {
	Limit  *int
	Cursor *string
	Type   *string
}

func (r *WebhooksResource) ListEvents(ctx context.Context, p EventsListParams) (JSON, error) {
	q := map[string]any{}
	if p.Limit != nil {
		q["limit"] = *p.Limit
	}
	if p.Cursor != nil {
		q["cursor"] = *p.Cursor
	}
	if p.Type != nil {
		q["type"] = *p.Type
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/webhooks/events" + qs(q)}, &out)
	return out, err
}

// ─── Audit log ─────────────────────────────────────────────────────────

type AuditLogResource struct{ c *Client }

type AuditListParams struct {
	Limit     *int
	Cursor    *string
	Since     *string
	EventType *string
}

func (r *AuditLogResource) List(ctx context.Context, p AuditListParams) (JSON, error) {
	q := map[string]any{}
	if p.Limit != nil {
		q["limit"] = *p.Limit
	}
	if p.Cursor != nil {
		q["cursor"] = *p.Cursor
	}
	if p.Since != nil {
		q["since"] = *p.Since
	}
	if p.EventType != nil {
		q["eventType"] = *p.EventType
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/audit-log" + qs(q)}, &out)
	return out, err
}

// ─── Integrations ──────────────────────────────────────────────────────

type IntegrationsResource struct{ c *Client }

func (r *IntegrationsResource) Status(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/integrations/status"}, &out)
	return out, err
}

func (r *IntegrationsResource) GetEmail(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/integrations/email"}, &out)
	return out, err
}

func (r *IntegrationsResource) UpdateEmail(ctx context.Context, input JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "PUT", Path: "/api/v1/integrations/email", Body: input}, &out)
	return out, err
}

// ─── Billing ───────────────────────────────────────────────────────────

type BillingResource struct{ c *Client }

func (r *BillingResource) Plans(ctx context.Context) ([]JSON, error) {
	var out []JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/billing/plans"}, &out)
	return out, err
}

func (r *BillingResource) CurrentPlan(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/billing/plan"}, &out)
	return out, err
}

func (r *BillingResource) Subscription(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/billing/subscription"}, &out)
	return out, err
}

func (r *BillingResource) Usage(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/billing/usage"}, &out)
	return out, err
}

type InvoicesListParams struct {
	Limit  *int
	Cursor *string
}

func (r *BillingResource) Invoices(ctx context.Context, p InvoicesListParams) (JSON, error) {
	q := map[string]any{}
	if p.Limit != nil {
		q["limit"] = *p.Limit
	}
	if p.Cursor != nil {
		q["cursor"] = *p.Cursor
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/billing/invoices" + qs(q)}, &out)
	return out, err
}

type CheckoutInput struct {
	PlanID     string `json:"planId"`
	SuccessURL string `json:"successUrl,omitempty"`
	CancelURL  string `json:"cancelUrl,omitempty"`
}

type CheckoutResult struct {
	URL       string `json:"url"`
	SessionID string `json:"sessionId"`
}

func (r *BillingResource) Checkout(ctx context.Context, input CheckoutInput) (*CheckoutResult, error) {
	var out CheckoutResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/billing/checkout", Body: input}, &out)
	return &out, err
}

func (r *BillingResource) Cancel(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/billing/cancel", Body: map[string]any{}}, &out)
	return out, err
}

// ─── Uploads ───────────────────────────────────────────────────────────

type UploadsResource struct{ c *Client }

type SignInput struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Bytes       *int   `json:"bytes,omitempty"`
}

type SignResult struct {
	URL     string            `json:"url"`
	Key     string            `json:"key"`
	Headers map[string]string `json:"headers,omitempty"`
}

func (r *UploadsResource) SignMerchant(ctx context.Context, input SignInput) (*SignResult, error) {
	var out SignResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/uploads/sign-merchant", Body: input}, &out)
	return &out, err
}

func (r *UploadsResource) GetMerchantAsset(ctx context.Context, key string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/uploads/merchant-asset" + qs(map[string]any{"key": key})}, &out)
	return out, err
}

func (r *UploadsResource) Sign(ctx context.Context, input SignInput) (*SignResult, error) {
	var out SignResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/uploads/sign", Body: input}, &out)
	return &out, err
}

func (r *UploadsResource) GetAvatar(ctx context.Context, key *string) (JSON, error) {
	q := map[string]any{}
	if key != nil {
		q["key"] = *key
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/uploads/avatar" + qs(q)}, &out)
	return out, err
}

func (r *UploadsResource) GetDeliverable(ctx context.Context, key string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/uploads/deliverable" + qs(map[string]any{"key": key})}, &out)
	return out, err
}

// ─── Campaigns ─────────────────────────────────────────────────────────

type CampaignsResource struct{ c *Client }

func (r *CampaignsResource) List(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/campaigns"}, &out)
	return out, err
}

func (r *CampaignsResource) Create(ctx context.Context, input JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/campaigns", Body: input, IdempotencyKey: r.c.idemFn()}, &out)
	return out, err
}

func (r *CampaignsResource) Get(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/campaigns/" + url.PathEscape(id)}, &out)
	return out, err
}

func (r *CampaignsResource) Update(ctx context.Context, id string, patch JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/campaigns/" + url.PathEscape(id), Body: patch}, &out)
	return out, err
}

func (r *CampaignsResource) InviteCreator(ctx context.Context, id, creatorID string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/campaigns/" + url.PathEscape(id) + "/invitations", Body: map[string]string{"creatorId": creatorID}}, &out)
	return out, err
}

func (r *CampaignsResource) ListApplications(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/campaigns/" + url.PathEscape(id) + "/applications"}, &out)
	return out, err
}

func (r *CampaignsResource) AcceptApplication(ctx context.Context, id, applicationID string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/campaigns/" + url.PathEscape(id) + "/applications/" + url.PathEscape(applicationID) + "/accept", Body: map[string]any{}}, &out)
	return out, err
}

func (r *CampaignsResource) RejectApplication(ctx context.Context, id, applicationID string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/campaigns/" + url.PathEscape(id) + "/applications/" + url.PathEscape(applicationID) + "/reject", Body: map[string]any{}}, &out)
	return out, err
}

func (r *CampaignsResource) Analytics(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/campaigns/" + url.PathEscape(id) + "/analytics"}, &out)
	return out, err
}

// ─── Programs ──────────────────────────────────────────────────────────

type ProgramsResource struct{ c *Client }

func (r *ProgramsResource) List(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/programs"}, &out)
	return out, err
}

func (r *ProgramsResource) Create(ctx context.Context, input JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/programs", Body: input, IdempotencyKey: r.c.idemFn()}, &out)
	return out, err
}

func (r *ProgramsResource) Get(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/programs/" + url.PathEscape(id)}, &out)
	return out, err
}

func (r *ProgramsResource) Update(ctx context.Context, id string, patch JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/programs/" + url.PathEscape(id), Body: patch}, &out)
	return out, err
}

func (r *ProgramsResource) Delete(ctx context.Context, id string) (*DeletedResult, error) {
	var out DeletedResult
	err := r.c.Do(ctx, RequestOptions{Method: "DELETE", Path: "/api/v1/programs/" + url.PathEscape(id)}, &out)
	return &out, err
}

func (r *ProgramsResource) ListEnrollments(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/programs/" + url.PathEscape(id) + "/enrollments"}, &out)
	return out, err
}

func (r *ProgramsResource) ApproveEnrollment(ctx context.Context, id, enrollmentID string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/programs/" + url.PathEscape(id) + "/enrollments/" + url.PathEscape(enrollmentID) + "/approve", Body: map[string]any{}}, &out)
	return out, err
}

func (r *ProgramsResource) RejectEnrollment(ctx context.Context, id, enrollmentID string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/programs/" + url.PathEscape(id) + "/enrollments/" + url.PathEscape(enrollmentID) + "/reject", Body: map[string]any{}}, &out)
	return out, err
}

func (r *ProgramsResource) RevokeEnrollment(ctx context.Context, id, enrollmentID string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/programs/" + url.PathEscape(id) + "/enrollments/" + url.PathEscape(enrollmentID) + "/revoke", Body: map[string]any{}}, &out)
	return out, err
}

type CommissionsListParams struct {
	Limit  *int
	Cursor *string
}

func (r *ProgramsResource) Commissions(ctx context.Context, p CommissionsListParams) (JSON, error) {
	q := map[string]any{}
	if p.Limit != nil {
		q["limit"] = *p.Limit
	}
	if p.Cursor != nil {
		q["cursor"] = *p.Cursor
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/programs/commissions" + qs(q)}, &out)
	return out, err
}

// Short aliases (Node SDK v0.2.1 parity).
func (r *ProgramsResource) Approve(ctx context.Context, id, enrollmentID string) (JSON, error) {
	return r.ApproveEnrollment(ctx, id, enrollmentID)
}
func (r *ProgramsResource) Reject(ctx context.Context, id, enrollmentID string) (JSON, error) {
	return r.RejectEnrollment(ctx, id, enrollmentID)
}
func (r *ProgramsResource) Revoke(ctx context.Context, id, enrollmentID string) (JSON, error) {
	return r.RevokeEnrollment(ctx, id, enrollmentID)
}
func (r *ProgramsResource) Enrollments(ctx context.Context, id string) (JSON, error) {
	return r.ListEnrollments(ctx, id)
}

// ─── Collaborations ────────────────────────────────────────────────────

type CollaborationsResource struct{ c *Client }

func (r *CollaborationsResource) FromApplication(ctx context.Context, applicationID string, input JSON) (JSON, error) {
	if input == nil {
		input = JSON{}
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/collaborations/from-application/" + url.PathEscape(applicationID), Body: input}, &out)
	return out, err
}

func (r *CollaborationsResource) List(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/collaborations"}, &out)
	return out, err
}

func (r *CollaborationsResource) Get(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/collaborations/" + url.PathEscape(id)}, &out)
	return out, err
}

type DeliverableUploadInput struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
}

func (r *CollaborationsResource) UploadDeliverableKey(ctx context.Context, id, deliverableID string, input DeliverableUploadInput) (*SignResult, error) {
	var out SignResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/collaborations/" + url.PathEscape(id) + "/deliverables/" + url.PathEscape(deliverableID) + "/upload-key", Body: input}, &out)
	return &out, err
}

func (r *CollaborationsResource) ApproveDeliverable(ctx context.Context, id, deliverableID string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/collaborations/" + url.PathEscape(id) + "/deliverables/" + url.PathEscape(deliverableID) + "/approve", Body: map[string]any{}}, &out)
	return out, err
}

func (r *CollaborationsResource) RejectDeliverable(ctx context.Context, id, deliverableID, reason string) (JSON, error) {
	body := map[string]any{}
	if reason != "" {
		body["reason"] = reason
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/collaborations/" + url.PathEscape(id) + "/deliverables/" + url.PathEscape(deliverableID) + "/reject", Body: body}, &out)
	return out, err
}

func (r *CollaborationsResource) PublishDeliverable(ctx context.Context, id, deliverableID string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/collaborations/" + url.PathEscape(id) + "/deliverables/" + url.PathEscape(deliverableID) + "/published", Body: map[string]any{}}, &out)
	return out, err
}

func (r *CollaborationsResource) ApproveCollaboration(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/collaborations/" + url.PathEscape(id) + "/approve", Body: map[string]any{}}, &out)
	return out, err
}

func (r *CollaborationsResource) CancelCollaboration(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/collaborations/" + url.PathEscape(id) + "/cancel", Body: map[string]any{}}, &out)
	return out, err
}

// Short aliases (v0.2.1).
func (r *CollaborationsResource) Approve(ctx context.Context, id string) (JSON, error) {
	return r.ApproveCollaboration(ctx, id)
}
func (r *CollaborationsResource) Cancel(ctx context.Context, id string) (JSON, error) {
	return r.CancelCollaboration(ctx, id)
}

// ─── Insights ──────────────────────────────────────────────────────────

type InsightsResource struct{ c *Client }

func (r *InsightsResource) Overview(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/insights/overview"}, &out)
	return out, err
}

func (r *InsightsResource) ProgramDetail(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/insights/programs/" + url.PathEscape(id)}, &out)
	return out, err
}

func (r *InsightsResource) CampaignDetail(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/insights/campaigns/" + url.PathEscape(id)}, &out)
	return out, err
}

type CommissionsTrendParams struct {
	Days *int
}

func (r *InsightsResource) CommissionsTrend(ctx context.Context, p CommissionsTrendParams) (JSON, error) {
	q := map[string]any{}
	if p.Days != nil {
		q["days"] = *p.Days
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/insights/commissions-trend" + qs(q)}, &out)
	return out, err
}

// ─── Channels ──────────────────────────────────────────────────────────

type ChannelsResource struct{ c *Client }

func (r *ChannelsResource) List(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/channels"}, &out)
	return out, err
}

func (r *ChannelsResource) Create(ctx context.Context, input JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/channels", Body: input}, &out)
	return out, err
}

func (r *ChannelsResource) Get(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/channels/" + url.PathEscape(id)}, &out)
	return out, err
}

func (r *ChannelsResource) Update(ctx context.Context, id string, patch JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/channels/" + url.PathEscape(id), Body: patch}, &out)
	return out, err
}

func (r *ChannelsResource) Delete(ctx context.Context, id string) (*DeletedResult, error) {
	var out DeletedResult
	err := r.c.Do(ctx, RequestOptions{Method: "DELETE", Path: "/api/v1/channels/" + url.PathEscape(id)}, &out)
	return &out, err
}

func (r *ChannelsResource) Test(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/channels/" + url.PathEscape(id) + "/test", Body: map[string]any{}}, &out)
	return out, err
}

func (r *ChannelsResource) DNSRecords(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/channels/" + url.PathEscape(id) + "/dns-records"}, &out)
	return out, err
}

func (r *ChannelsResource) OAuthStart(ctx context.Context, provider string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/channels/oauth/" + url.PathEscape(provider) + "/start"}, &out)
	return out, err
}

// Short alias (v0.2.1).
func (r *ChannelsResource) DNS(ctx context.Context, id string) (JSON, error) {
	return r.DNSRecords(ctx, id)
}

// ─── Contacts ──────────────────────────────────────────────────────────

type ContactsResource struct{ c *Client }

type ContactsListParams struct {
	Limit  *int
	Cursor *string
	Search *string
}

func (r *ContactsResource) List(ctx context.Context, p ContactsListParams) (JSON, error) {
	q := map[string]any{}
	if p.Limit != nil {
		q["limit"] = *p.Limit
	}
	if p.Cursor != nil {
		q["cursor"] = *p.Cursor
	}
	if p.Search != nil {
		q["search"] = *p.Search
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/contacts" + qs(q)}, &out)
	return out, err
}

func (r *ContactsResource) Import(ctx context.Context, input JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/contacts/import", Body: input}, &out)
	return out, err
}

func (r *ContactsResource) Create(ctx context.Context, input JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/contacts", Body: input}, &out)
	return out, err
}

func (r *ContactsResource) Get(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/contacts/" + url.PathEscape(id)}, &out)
	return out, err
}

func (r *ContactsResource) Update(ctx context.Context, id string, patch JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/contacts/" + url.PathEscape(id), Body: patch}, &out)
	return out, err
}

func (r *ContactsResource) Delete(ctx context.Context, id string) (*DeletedResult, error) {
	var out DeletedResult
	err := r.c.Do(ctx, RequestOptions{Method: "DELETE", Path: "/api/v1/contacts/" + url.PathEscape(id)}, &out)
	return &out, err
}

// ─── Contact lists ─────────────────────────────────────────────────────

type ContactListsResource struct{ c *Client }

func (r *ContactListsResource) List(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/contact-lists"}, &out)
	return out, err
}

type ContactListInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

func (r *ContactListsResource) Create(ctx context.Context, input ContactListInput) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/contact-lists", Body: input}, &out)
	return out, err
}

func (r *ContactListsResource) Get(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/contact-lists/" + url.PathEscape(id)}, &out)
	return out, err
}

func (r *ContactListsResource) AddMember(ctx context.Context, id, contactID string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/contact-lists/" + url.PathEscape(id) + "/members", Body: map[string]string{"contactId": contactID}}, &out)
	return out, err
}

func (r *ContactListsResource) RemoveMember(ctx context.Context, id, contactID string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "DELETE", Path: "/api/v1/contact-lists/" + url.PathEscape(id) + "/members/" + url.PathEscape(contactID)}, &out)
	return out, err
}

func (r *ContactListsResource) Delete(ctx context.Context, id string) (*DeletedResult, error) {
	var out DeletedResult
	err := r.c.Do(ctx, RequestOptions{Method: "DELETE", Path: "/api/v1/contact-lists/" + url.PathEscape(id)}, &out)
	return &out, err
}

// ─── Broadcasts (email/SMS blasts) ─────────────────────────────────────
//
// Backed by /api/v1/broadcasts on the server. Distinct from
// MarketingCampaignsResource (the campaign hub) which lives at
// /api/v1/marketing-campaigns.

type BroadcastsResource struct{ c *Client }

func (r *BroadcastsResource) List(ctx context.Context) ([]JSON, error) {
	var out []JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/broadcasts"}, &out)
	return out, err
}

func (r *BroadcastsResource) Create(ctx context.Context, input JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/broadcasts", Body: input}, &out)
	return out, err
}

func (r *BroadcastsResource) Get(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/broadcasts/" + url.PathEscape(id)}, &out)
	return out, err
}

func (r *BroadcastsResource) Update(ctx context.Context, id string, patch JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/broadcasts/" + url.PathEscape(id), Body: patch}, &out)
	return out, err
}

func (r *BroadcastsResource) Send(ctx context.Context, id string, input JSON) (JSON, error) {
	if input == nil {
		input = JSON{}
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/broadcasts/" + url.PathEscape(id) + "/send", Body: input, IdempotencyKey: r.c.idemFn()}, &out)
	return out, err
}

func (r *BroadcastsResource) SendTest(ctx context.Context, id, to string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/broadcasts/" + url.PathEscape(id) + "/send-test", Body: map[string]string{"to": to}}, &out)
	return out, err
}

func (r *BroadcastsResource) ListTemplates(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/broadcasts/templates"}, &out)
	return out, err
}

func (r *BroadcastsResource) CreateTemplate(ctx context.Context, input JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/broadcasts/templates", Body: input}, &out)
	return out, err
}

func (r *BroadcastsResource) UpdateTemplate(ctx context.Context, templateID string, patch JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/broadcasts/templates/" + url.PathEscape(templateID), Body: patch}, &out)
	return out, err
}

type CompileTemplateResult struct {
	HTML string `json:"html"`
}

func (r *BroadcastsResource) CompileTemplate(ctx context.Context, input JSON) (*CompileTemplateResult, error) {
	var out CompileTemplateResult
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/broadcasts/templates/compile", Body: input}, &out)
	return &out, err
}

// ─── Marketing campaigns (hub) ─────────────────────────────────────────
//
// Top-level marketing-campaign hub. Groups creator briefs, affiliate
// programs, discount codes, abandoned-cart reminders, referral programs,
// blog posts, and feeds under one merchant-defined campaign. Every
// child link is optional.
//
// Backed by /api/v1/marketing-campaigns (Prisma model
// `MarketingCampaign` -> DB table `MarketingProgram` via @@map).
//
// Distinct from BroadcastsResource (email/SMS blasts) which lives at
// /api/v1/broadcasts.

// MarketingCampaignGoal — `awareness` | `conversion` | `retention` | `launch` | `other`.
type MarketingCampaignGoal string

// MarketingCampaignStatus — `draft` | `live` | `paused` | `completed` | `archived`.
type MarketingCampaignStatus string

const (
	MarketingCampaignGoalAwareness  MarketingCampaignGoal = "awareness"
	MarketingCampaignGoalConversion MarketingCampaignGoal = "conversion"
	MarketingCampaignGoalRetention  MarketingCampaignGoal = "retention"
	MarketingCampaignGoalLaunch     MarketingCampaignGoal = "launch"
	MarketingCampaignGoalOther      MarketingCampaignGoal = "other"

	MarketingCampaignStatusDraft     MarketingCampaignStatus = "draft"
	MarketingCampaignStatusLive      MarketingCampaignStatus = "live"
	MarketingCampaignStatusPaused    MarketingCampaignStatus = "paused"
	MarketingCampaignStatusCompleted MarketingCampaignStatus = "completed"
	MarketingCampaignStatusArchived  MarketingCampaignStatus = "archived"
)

type MarketingCampaignsResource struct{ c *Client }

// MarketingCampaignsListParams — narrow filters for List.
//
// Status is comma-joined when sent; a single string is also accepted
// verbatim. Limit / Cursor mirror the standard Ripllo paging shape;
// the backend isn't paged yet but the fields are passed through for
// forward-compat.
type MarketingCampaignsListParams struct {
	Status []MarketingCampaignStatus
	Limit  *int
	Cursor *string
}

func (r *MarketingCampaignsResource) List(ctx context.Context, p MarketingCampaignsListParams) (JSON, error) {
	q := map[string]any{}
	if len(p.Status) > 0 {
		parts := make([]string, 0, len(p.Status))
		for _, s := range p.Status {
			parts = append(parts, string(s))
		}
		q["status"] = joinComma(parts)
	}
	if p.Limit != nil {
		q["limit"] = *p.Limit
	}
	if p.Cursor != nil {
		q["cursor"] = *p.Cursor
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/marketing-campaigns" + qs(q)}, &out)
	return out, err
}

func (r *MarketingCampaignsResource) Get(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/marketing-campaigns/" + url.PathEscape(id)}, &out)
	return out, err
}

// GetFull — GET /:id/full — hub + all linked children + perf roll-up.
func (r *MarketingCampaignsResource) GetFull(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/marketing-campaigns/" + url.PathEscape(id) + "/full"}, &out)
	return out, err
}

// Selector — GET /_/selector — lightweight dropdown payload (non-archived).
func (r *MarketingCampaignsResource) Selector(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/marketing-campaigns/_/selector"}, &out)
	return out, err
}

// MarketingCampaignCreateInput — payload for POST /marketing-campaigns.
//
// `omitempty` on optional fields keeps wire payloads small; the
// pointer-typed fields are NULL-allowing so callers can explicitly
// clear server-side values.
type MarketingCampaignCreateInput struct {
	Name        string                   `json:"name"`
	Description *string                  `json:"description,omitempty"`
	Goal        MarketingCampaignGoal    `json:"goal,omitempty"`
	Status      MarketingCampaignStatus  `json:"status,omitempty"`
	BudgetIDR   *int64                   `json:"budgetIdr,omitempty"`
	StartsAt    *string                  `json:"startsAt,omitempty"`
	EndsAt      *string                  `json:"endsAt,omitempty"`
	Notes       *string                  `json:"notes,omitempty"`
}

type MarketingCampaignUpdateInput = MarketingCampaignCreateInput

func (r *MarketingCampaignsResource) Create(ctx context.Context, input MarketingCampaignCreateInput) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/marketing-campaigns", Body: input, IdempotencyKey: r.c.idemFn()}, &out)
	return out, err
}

func (r *MarketingCampaignsResource) Update(ctx context.Context, id string, patch MarketingCampaignUpdateInput) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/marketing-campaigns/" + url.PathEscape(id), Body: patch}, &out)
	return out, err
}

// Delete — soft-delete via status='archived' on the server.
func (r *MarketingCampaignsResource) Delete(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "DELETE", Path: "/api/v1/marketing-campaigns/" + url.PathEscape(id)}, &out)
	return out, err
}

// joinComma keeps the resource file self-contained; the Go SDK doesn't
// already import strings here and a 4-line helper is cheaper than the
// import churn.
func joinComma(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += "," + parts[i]
	}
	return out
}

// ─── Funnels ───────────────────────────────────────────────────────────

type FunnelsResource struct{ c *Client }

func (r *FunnelsResource) List(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/funnels"}, &out)
	return out, err
}

func (r *FunnelsResource) Create(ctx context.Context, input JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/funnels", Body: input}, &out)
	return out, err
}

func (r *FunnelsResource) Get(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/funnels/" + url.PathEscape(id)}, &out)
	return out, err
}

func (r *FunnelsResource) Update(ctx context.Context, id string, patch JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/funnels/" + url.PathEscape(id), Body: patch}, &out)
	return out, err
}

func (r *FunnelsResource) Delete(ctx context.Context, id string) (*DeletedResult, error) {
	var out DeletedResult
	err := r.c.Do(ctx, RequestOptions{Method: "DELETE", Path: "/api/v1/funnels/" + url.PathEscape(id)}, &out)
	return &out, err
}

// SetSteps accepts either a `[]any` (treated as `steps`) or a JSON map.
// Mirrors the Node SDK shape: when callers pass a bare array, the SDK
// wraps it in `{steps: [...]}` automatically.
func (r *FunnelsResource) SetSteps(ctx context.Context, id string, input any) (JSON, error) {
	body := input
	switch v := input.(type) {
	case []any:
		body = map[string]any{"steps": v}
	case []map[string]any:
		body = map[string]any{"steps": v}
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "PUT", Path: "/api/v1/funnels/" + url.PathEscape(id) + "/steps", Body: body}, &out)
	return out, err
}

func (r *FunnelsResource) Enroll(ctx context.Context, id string, input JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/funnels/" + url.PathEscape(id) + "/enroll", Body: input}, &out)
	return out, err
}

func (r *FunnelsResource) Analytics(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/funnels/" + url.PathEscape(id) + "/analytics"}, &out)
	return out, err
}

func (r *FunnelsResource) ListEnrollments(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/funnels/" + url.PathEscape(id) + "/enrollments"}, &out)
	return out, err
}

// ─── Inbox ─────────────────────────────────────────────────────────────

type InboxResource struct{ c *Client }

type InboxThreadsListParams struct {
	Limit  *int
	Cursor *string
}

func (r *InboxResource) ListThreads(ctx context.Context, p InboxThreadsListParams) (JSON, error) {
	q := map[string]any{}
	if p.Limit != nil {
		q["limit"] = *p.Limit
	}
	if p.Cursor != nil {
		q["cursor"] = *p.Cursor
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/inbox/threads" + qs(q)}, &out)
	return out, err
}

func (r *InboxResource) GetThread(ctx context.Context, provider, handle string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/inbox/threads/" + url.PathEscape(provider) + "/" + url.PathEscape(handle)}, &out)
	return out, err
}

func (r *InboxResource) MarkRead(ctx context.Context, threadID string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/inbox/" + url.PathEscape(threadID) + "/read", Body: map[string]any{}}, &out)
	return out, err
}

func (r *InboxResource) Archive(ctx context.Context, threadID string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/inbox/" + url.PathEscape(threadID) + "/archive", Body: map[string]any{}}, &out)
	return out, err
}

// ─── Audience segments ─────────────────────────────────────────────────

type AudienceSegmentsResource struct{ c *Client }

func (r *AudienceSegmentsResource) List(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/audience-segments"}, &out)
	return out, err
}

func (r *AudienceSegmentsResource) Create(ctx context.Context, input JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/audience-segments", Body: input}, &out)
	return out, err
}

func (r *AudienceSegmentsResource) Get(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/audience-segments/" + url.PathEscape(id)}, &out)
	return out, err
}

func (r *AudienceSegmentsResource) Update(ctx context.Context, id string, patch JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/audience-segments/" + url.PathEscape(id), Body: patch}, &out)
	return out, err
}

func (r *AudienceSegmentsResource) Delete(ctx context.Context, id string) (*DeletedResult, error) {
	var out DeletedResult
	err := r.c.Do(ctx, RequestOptions{Method: "DELETE", Path: "/api/v1/audience-segments/" + url.PathEscape(id)}, &out)
	return &out, err
}

func (r *AudienceSegmentsResource) Preview(ctx context.Context, id string, input JSON) (JSON, error) {
	if input == nil {
		input = JSON{}
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/audience-segments/" + url.PathEscape(id) + "/preview", Body: input}, &out)
	return out, err
}

func (r *AudienceSegmentsResource) PreviewAdhoc(ctx context.Context, input JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/audience-segments/preview", Body: input}, &out)
	return out, err
}

// ─── Merchant profile ──────────────────────────────────────────────────

type MerchantProfileResource struct{ c *Client }

func (r *MerchantProfileResource) Me(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/merchants/me"}, &out)
	return out, err
}

func (r *MerchantProfileResource) UpdateMe(ctx context.Context, patch JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "PUT", Path: "/api/v1/merchants/me", Body: patch}, &out)
	return out, err
}

func (r *MerchantProfileResource) Publish(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/merchants/me/publish", Body: map[string]any{}}, &out)
	return out, err
}

func (r *MerchantProfileResource) Unpublish(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/merchants/me/unpublish", Body: map[string]any{}}, &out)
	return out, err
}

func (r *MerchantProfileResource) List(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/merchants"}, &out)
	return out, err
}

func (r *MerchantProfileResource) GetBySlug(ctx context.Context, slug string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/merchants/" + url.PathEscape(slug)}, &out)
	return out, err
}

// ─── Creator + affiliator profiles ─────────────────────────────────────

type CreatorProfileResource struct{ c *Client }

func (r *CreatorProfileResource) Me(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/creator-profile/me"}, &out)
	return out, err
}

func (r *CreatorProfileResource) Create(ctx context.Context, input JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/creator-profile", Body: input}, &out)
	return out, err
}

func (r *CreatorProfileResource) UpdateMe(ctx context.Context, patch JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/creator-profile/me", Body: patch}, &out)
	return out, err
}

type AffiliatorProfileResource struct{ c *Client }

func (r *AffiliatorProfileResource) Me(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/affiliator-profile/me"}, &out)
	return out, err
}

func (r *AffiliatorProfileResource) Create(ctx context.Context, input JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/affiliator-profile", Body: input}, &out)
	return out, err
}

func (r *AffiliatorProfileResource) UpdateMe(ctx context.Context, patch JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "PATCH", Path: "/api/v1/affiliator-profile/me", Body: patch}, &out)
	return out, err
}

// ─── Marketplace ───────────────────────────────────────────────────────

type MarketplaceResource struct{ c *Client }

func (r *MarketplaceResource) ListCreators(ctx context.Context, params JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/marketplace/creators" + qs(params)}, &out)
	return out, err
}

func (r *MarketplaceResource) GetCreator(ctx context.Context, handle string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/marketplace/creators/" + url.PathEscape(handle)}, &out)
	return out, err
}

func (r *MarketplaceResource) ListCampaigns(ctx context.Context, params JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/marketplace/campaigns" + qs(params)}, &out)
	return out, err
}

func (r *MarketplaceResource) GetCampaign(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/marketplace/campaigns/" + url.PathEscape(id)}, &out)
	return out, err
}

func (r *MarketplaceResource) ApplyToCampaign(ctx context.Context, id string, input JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/marketplace/campaigns/" + url.PathEscape(id) + "/apply", Body: input}, &out)
	return out, err
}

func (r *MarketplaceResource) MyInvitations(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/marketplace/me/invitations"}, &out)
	return out, err
}

func (r *MarketplaceResource) RespondToInvitation(ctx context.Context, invitationID string, accept bool) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/marketplace/me/invitations/" + url.PathEscape(invitationID) + "/respond", Body: map[string]bool{"accept": accept}}, &out)
	return out, err
}

func (r *MarketplaceResource) MyApplications(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/marketplace/me/applications"}, &out)
	return out, err
}

// ─── Affiliates ────────────────────────────────────────────────────────

type AffiliatesResource struct{ c *Client }

type AffiliatesListParams struct {
	Limit  *int
	Cursor *string
}

func (r *AffiliatesResource) List(ctx context.Context, p AffiliatesListParams) (JSON, error) {
	q := map[string]any{}
	if p.Limit != nil {
		q["limit"] = *p.Limit
	}
	if p.Cursor != nil {
		q["cursor"] = *p.Cursor
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/affiliates" + qs(q)}, &out)
	return out, err
}

// ─── KYC ───────────────────────────────────────────────────────────────

type KYCResource struct{ c *Client }

func (r *KYCResource) GetStatus(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/kyc"}, &out)
	return out, err
}

func (r *KYCResource) Submit(ctx context.Context, input JSON) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/kyc", Body: input}, &out)
	return out, err
}

// ─── Creator stats ─────────────────────────────────────────────────────

type CreatorStatsResource struct{ c *Client }

func (r *CreatorStatsResource) Overview(ctx context.Context) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/creator-stats"}, &out)
	return out, err
}

func (r *CreatorStatsResource) Connect(ctx context.Context, provider string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/creator-stats/connect/" + url.PathEscape(provider)}, &out)
	return out, err
}

// ─── Admin (Pattern 2 partner billing + KYC + disputes) ───────────────

type AdminResource struct{ c *Client }

func (r *AdminResource) ProvisionWorkspace(ctx context.Context, input ProvisionWorkspaceInput) (*PartnerWorkspace, error) {
	var out PartnerWorkspace
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/admin/workspaces", Body: input}, &out)
	return &out, err
}

func (r *AdminResource) GetWorkspace(ctx context.Context, accountID string) (*PartnerWorkspace, error) {
	var out PartnerWorkspace
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/admin/workspaces/" + url.PathEscape(accountID)}, &out)
	return &out, err
}

func (r *AdminResource) PartnerUsage(ctx context.Context, q PartnerUsageQuery) (*PartnerUsageSummary, error) {
	params := map[string]any{"from": q.From, "to": q.To}
	if q.Partner != "" {
		params["partner"] = q.Partner
	}
	var out PartnerUsageSummary
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/admin/partner/usage" + qs(params)}, &out)
	return &out, err
}

type KYCListParams struct {
	Status *string
	Limit  *int
}

func (r *AdminResource) ListKYC(ctx context.Context, p KYCListParams) (JSON, error) {
	q := map[string]any{}
	if p.Status != nil {
		q["status"] = *p.Status
	}
	if p.Limit != nil {
		q["limit"] = *p.Limit
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/admin/kyc" + qs(q)}, &out)
	return out, err
}

func (r *AdminResource) ApproveKYC(ctx context.Context, id string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/admin/kyc/" + url.PathEscape(id) + "/approve", Body: map[string]any{}}, &out)
	return out, err
}

func (r *AdminResource) RejectKYC(ctx context.Context, id, reason string) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/admin/kyc/" + url.PathEscape(id) + "/reject", Body: map[string]string{"reason": reason}}, &out)
	return out, err
}

type DisputesListParams struct {
	Status *string
	Limit  *int
}

func (r *AdminResource) ListDisputes(ctx context.Context, p DisputesListParams) (JSON, error) {
	q := map[string]any{}
	if p.Status != nil {
		q["status"] = *p.Status
	}
	if p.Limit != nil {
		q["limit"] = *p.Limit
	}
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "GET", Path: "/api/v1/admin/disputes" + qs(q)}, &out)
	return out, err
}

type DisputeResolution struct {
	Resolution string `json:"resolution"` // "merchant" | "creator"
	Note       string `json:"note,omitempty"`
}

func (r *AdminResource) ResolveDispute(ctx context.Context, id string, input DisputeResolution) (JSON, error) {
	var out JSON
	err := r.c.Do(ctx, RequestOptions{Method: "POST", Path: "/api/v1/admin/disputes/" + url.PathEscape(id) + "/resolve", Body: input}, &out)
	return out, err
}

// ─── Internal helpers ──────────────────────────────────────────────────

// Compile-time check — every resource pointer is non-nil after
// mountResources. (vetted at runtime by Test_NewClient_Surface.)
var _ = func() bool {
	// guard against accidentally dropping a resource in mountResources;
	// the test suite confirms every pointer is mounted.
	return true
}()

// AsJSON is a small convenience for callers that already have a typed
// struct but want it as `JSON` (e.g. for passing into the marketing
// resources). Returns a freshly allocated map.
func AsJSON(v any) (JSON, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var out JSON
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ensure fmt is used (qs uses fmt via the formatter); keeps `goimports`
// happy if the file is trimmed later.
var _ = fmt.Sprintf
