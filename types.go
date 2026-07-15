package ripllo

import "encoding/json"

// ─── Envelope ──────────────────────────────────────────────────────────

// apiEnvelope wraps every Ripllo HTTP response. Either Data or Error
// is non-nil. Meta is best-effort.
type apiEnvelope struct {
	Data  json.RawMessage `json:"data"`
	Error *envelopeError  `json:"error"`
	Meta  *envelopeMeta   `json:"meta,omitempty"`
}

type envelopeError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type envelopeMeta struct {
	RequestID string `json:"requestId,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Cursor    string `json:"cursor,omitempty"`
	HasMore   bool   `json:"hasMore,omitempty"`
}

// ─── Discount codes ────────────────────────────────────────────────────

// DiscountType — `percent` | `fixed` | `shipping_percent` | `shipping_fixed`.
type DiscountType string

// DiscountScope — `cart` | `products` | `tags`.
type DiscountScope string

const (
	DiscountTypePercent         DiscountType = "percent"
	DiscountTypeFixed           DiscountType = "fixed"
	DiscountTypeShippingPercent DiscountType = "shipping_percent"
	DiscountTypeShippingFixed   DiscountType = "shipping_fixed"

	DiscountScopeCart     DiscountScope = "cart"
	DiscountScopeProducts DiscountScope = "products"
	DiscountScopeTags     DiscountScope = "tags"
)

type DiscountCode struct {
	ID                 string        `json:"id"`
	AccountID          string        `json:"accountId"`
	Code               string        `json:"code"`
	Description        *string       `json:"description"`
	Type               DiscountType  `json:"type"`
	Value              float64       `json:"value"`
	Currency           string        `json:"currency"`
	Scope              DiscountScope `json:"scope"`
	ProductIDs         []string      `json:"productIds"`
	TagFilter          []string      `json:"tagFilter"`
	MinPurchaseAmount  *float64      `json:"minPurchaseAmount"`
	MaxUsesTotal       *int          `json:"maxUsesTotal"`
	MaxUsesPerCustomer *int          `json:"maxUsesPerCustomer"`
	StartsAt           *string       `json:"startsAt"`
	ExpiresAt          *string       `json:"expiresAt"`
	Active             bool          `json:"active"`
	Public             bool          `json:"public"`
	RedemptionCount    int           `json:"redemptionCount"`
	Source             *string       `json:"source"`
	SourceRefID        *string       `json:"sourceRefId"`
	CreatedAt          string        `json:"createdAt"`
	UpdatedAt          string        `json:"updatedAt"`
}

type DiscountCodeListPage struct {
	Items      []DiscountCode `json:"items"`
	Total      int            `json:"total"`
	NextCursor *string        `json:"nextCursor"`
	HasMore    bool           `json:"hasMore"`
}

// DiscountCodeCreateInput — fields for `discountCodes.Create`. Pointer
// fields are optional (omitted when nil).
type DiscountCodeCreateInput struct {
	Code               string        `json:"code"`
	Description        *string       `json:"description,omitempty"`
	Type               DiscountType  `json:"type"`
	Value              float64       `json:"value"`
	Currency           string        `json:"currency"`
	Scope              DiscountScope `json:"scope,omitempty"`
	ProductIDs         []string      `json:"productIds,omitempty"`
	TagFilter          []string      `json:"tagFilter,omitempty"`
	MinPurchaseAmount  *float64      `json:"minPurchaseAmount,omitempty"`
	MaxUsesTotal       *int          `json:"maxUsesTotal,omitempty"`
	MaxUsesPerCustomer *int          `json:"maxUsesPerCustomer,omitempty"`
	StartsAt           *string       `json:"startsAt,omitempty"`
	ExpiresAt          *string       `json:"expiresAt,omitempty"`
	Active             *bool         `json:"active,omitempty"`
	Public             *bool         `json:"public,omitempty"`
}

type ValidateLineItem struct {
	ProductID *string  `json:"productId,omitempty"`
	Price     float64  `json:"price"`
	Quantity  int      `json:"quantity"`
	Tags      []string `json:"tags,omitempty"`
}

type ValidateInput struct {
	AccountID    string             `json:"accountId"`
	Code         string             `json:"code"`
	Subtotal     float64            `json:"subtotal"`
	Currency     string             `json:"currency"`
	ShippingCost *float64           `json:"shippingCost,omitempty"`
	CustomerID   *string            `json:"customerId,omitempty"`
	Items        []ValidateLineItem `json:"items,omitempty"`
}

type ValidateResult struct {
	Valid             bool          `json:"valid"`
	Code              *DiscountCode `json:"code,omitempty"`
	Reason            string        `json:"reason,omitempty"`
	DiscountAmount    float64       `json:"discountAmount"`
	DiscountShipping  float64       `json:"discountShipping"`
}

type RedeemInput struct {
	AccountID         string   `json:"accountId"`
	DiscountCodeID    string   `json:"discountCodeId"`
	CheckoutSessionID string   `json:"checkoutSessionId"`
	CustomerID        *string  `json:"customerId,omitempty"`
	AppliedAmount     float64  `json:"appliedAmount"`
	AppliedShipping   *float64 `json:"appliedShipping,omitempty"`
	ExternalSource    *string  `json:"externalSource,omitempty"`
	ExternalRef       *string  `json:"externalRef,omitempty"`
}

type RedeemResult struct {
	ID      string `json:"id"`
	Created bool   `json:"created"`
}

type PublicApplicableCode struct {
	Code              string        `json:"code"`
	Description       *string       `json:"description"`
	Type              DiscountType  `json:"type"`
	Value             float64       `json:"value"`
	Scope             DiscountScope `json:"scope"`
	MinPurchaseAmount *float64      `json:"minPurchaseAmount"`
	ExpiresAt         *string       `json:"expiresAt"`
	EstimatedDiscount float64       `json:"estimatedDiscount"`
	Eligible          bool          `json:"eligible"`
}

type ArchiveResult struct {
	ID     string `json:"id"`
	Active bool   `json:"active"`
}

// ─── Pixels ────────────────────────────────────────────────────────────

type MerchantPixels struct {
	ID                     string  `json:"id,omitempty"`
	AccountID              string  `json:"accountId,omitempty"`
	MetaPixelID            *string `json:"metaPixelId"`
	MetaCapiAccessToken    *string `json:"metaCapiAccessToken,omitempty"`
	MetaTestEventCode      *string `json:"metaTestEventCode,omitempty"`
	GoogleAnalyticsID      *string `json:"googleAnalyticsId"`
	GoogleAdsConversionID  *string `json:"googleAdsConversionId"`
	GoogleAdsPurchaseLabel *string `json:"googleAdsPurchaseLabel"`
	TiktokPixelID          *string `json:"tiktokPixelId"`
	Enabled                bool    `json:"enabled"`
}

type PublicPixels struct {
	MetaPixelID            *string `json:"metaPixelId"`
	GoogleAnalyticsID      *string `json:"googleAnalyticsId"`
	GoogleAdsConversionID  *string `json:"googleAdsConversionId"`
	GoogleAdsPurchaseLabel *string `json:"googleAdsPurchaseLabel"`
	TiktokPixelID          *string `json:"tiktokPixelId"`
	Enabled                bool    `json:"enabled"`
}

// ─── Feeds ─────────────────────────────────────────────────────────────

type MerchantFeedConfig struct {
	Enabled                      bool    `json:"enabled"`
	DefaultGoogleProductCategory *string `json:"defaultGoogleProductCategory"`
	IncludeUnpublished           bool    `json:"includeUnpublished"`
}

// ─── Blog ──────────────────────────────────────────────────────────────

// BlogPostStatus — `draft` | `published`.
type BlogPostStatus string

const (
	BlogPostStatusDraft     BlogPostStatus = "draft"
	BlogPostStatusPublished BlogPostStatus = "published"
)

type BlogPost struct {
	ID              string         `json:"id"`
	AccountID       string         `json:"accountId"`
	Slug            string         `json:"slug"`
	Title           string         `json:"title"`
	Excerpt         *string        `json:"excerpt"`
	Body            string         `json:"body"`
	CoverImage      *string        `json:"coverImage"`
	Status          BlogPostStatus `json:"status"`
	PublishedAt     *string        `json:"publishedAt"`
	AuthorName      *string        `json:"authorName"`
	Tags            []string       `json:"tags"`
	MetaTitle       *string        `json:"metaTitle"`
	MetaDescription *string        `json:"metaDescription"`
	CreatedAt       string         `json:"createdAt"`
	UpdatedAt       string         `json:"updatedAt"`
}

type BlogPostInput struct {
	Slug            string          `json:"slug"`
	Title           string          `json:"title"`
	Excerpt         *string         `json:"excerpt,omitempty"`
	Body            string          `json:"body"`
	CoverImage      *string         `json:"coverImage,omitempty"`
	Status          *BlogPostStatus `json:"status,omitempty"`
	PublishedAt     *string         `json:"publishedAt,omitempty"`
	AuthorName      *string         `json:"authorName,omitempty"`
	Tags            []string        `json:"tags,omitempty"`
	MetaTitle       *string         `json:"metaTitle,omitempty"`
	MetaDescription *string         `json:"metaDescription,omitempty"`
}

type BlogPostWrap struct {
	Post BlogPost `json:"post"`
}

type BlogPostListWrap struct {
	Posts []BlogPost `json:"posts"`
}

// ─── Abandoned cart ────────────────────────────────────────────────────

type AbandonedCartConfig struct {
	ID             string  `json:"id,omitempty"`
	AccountID      string  `json:"accountId,omitempty"`
	Enabled        bool    `json:"enabled"`
	DelayHours     int     `json:"delayHours"`
	EmailSubject   string  `json:"emailSubject"`
	EmailPreview   string  `json:"emailPreview"`
	DiscountCodeID *string `json:"discountCodeId"`
}

type AbandonedCartReminder struct {
	ID                   string          `json:"id"`
	AccountID            string          `json:"accountId"`
	CustomerID           string          `json:"customerId"`
	CartID               string          `json:"cartId"`
	Email                string          `json:"email"`
	CartSnapshot         json.RawMessage `json:"cartSnapshot"`
	ValueAtSend          float64         `json:"valueAtSend"`
	CurrencyAtSend       string          `json:"currencyAtSend"`
	DiscountCodeID       *string         `json:"discountCodeId"`
	ExternalSource       *string         `json:"externalSource"`
	ExternalRef          *string         `json:"externalRef"`
	SentAt               string          `json:"sentAt"`
	RecoveredAt          *string         `json:"recoveredAt"`
	RecoveredBySessionID *string         `json:"recoveredBySessionId"`
}

type RecoveryStats struct {
	RemindersSent          int      `json:"remindersSent"`
	CartsRecovered         int      `json:"cartsRecovered"`
	RecoveryRate           float64  `json:"recoveryRate"`
	RecoveredValueAtSend   float64  `json:"recoveredValueAtSend"`
	Currency               *string  `json:"currency"`
}

type RecordReminderInput struct {
	AccountID      string          `json:"accountId"`
	CustomerID     string          `json:"customerId"`
	CartID         string          `json:"cartId"`
	Email          string          `json:"email"`
	CartSnapshot   any             `json:"cartSnapshot"`
	ValueAtSend    float64         `json:"valueAtSend"`
	CurrencyAtSend string          `json:"currencyAtSend"`
	DiscountCodeID *string         `json:"discountCodeId,omitempty"`
	ExternalSource *string         `json:"externalSource,omitempty"`
	ExternalRef    *string         `json:"externalRef,omitempty"`
}

type RecordReminderResult struct {
	ID      string `json:"id"`
	Created bool   `json:"created"`
	Reason  string `json:"reason,omitempty"`
}

type MarkRecoveredInput struct {
	AccountID         string  `json:"accountId"`
	CustomerID        string  `json:"customerId"`
	CheckoutSessionID string  `json:"checkoutSessionId"`
	CompletedAt       *string `json:"completedAt,omitempty"`
}

type MarkRecoveredResult struct {
	Recovered  bool   `json:"recovered"`
	ReminderID string `json:"reminderId,omitempty"`
}

type RemindersListWrap struct {
	Items []AbandonedCartReminder `json:"items"`
}

// ─── Referrals ─────────────────────────────────────────────────────────

// ReferralAttributionStatus — `pending` | `rewarded` | `voided` | `expired`.
type ReferralAttributionStatus string

type ReferralProgram struct {
	ID                    string       `json:"id"`
	AccountID             string       `json:"accountId"`
	Enabled               bool         `json:"enabled"`
	RewardType            DiscountType `json:"rewardType"`
	ReferrerValue         float64      `json:"referrerValue"`
	RefereeValue          float64      `json:"refereeValue"`
	Currency              string       `json:"currency"`
	MinPurchaseAmount     *float64     `json:"minPurchaseAmount"`
	RewardExpiryDays      int          `json:"rewardExpiryDays"`
	AttributionWindowDays int          `json:"attributionWindowDays"`
	MaxRewardsPerReferrer *int         `json:"maxRewardsPerReferrer"`
	ProgramTerms          *string      `json:"programTerms"`
	CreatedAt             string       `json:"createdAt"`
	UpdatedAt             string       `json:"updatedAt"`
}

type ReferralProgramInput struct {
	Enabled               *bool        `json:"enabled,omitempty"`
	RewardType            DiscountType `json:"rewardType"`
	ReferrerValue         float64      `json:"referrerValue"`
	RefereeValue          float64      `json:"refereeValue"`
	Currency              string       `json:"currency"`
	MinPurchaseAmount     *float64     `json:"minPurchaseAmount,omitempty"`
	RewardExpiryDays      *int         `json:"rewardExpiryDays,omitempty"`
	AttributionWindowDays *int         `json:"attributionWindowDays,omitempty"`
	MaxRewardsPerReferrer *int         `json:"maxRewardsPerReferrer,omitempty"`
	ProgramTerms          *string      `json:"programTerms,omitempty"`
}

type ReferralLink struct {
	ID         string  `json:"id"`
	ProgramID  string  `json:"programId"`
	AccountID  string  `json:"accountId"`
	CustomerID string  `json:"customerId"`
	Code       string  `json:"code"`
	Clicks     int     `json:"clicks"`
	Signups    int     `json:"signups"`
	Rewards    int     `json:"rewards"`
	Revenue    float64 `json:"revenue"`
	CreatedAt  string  `json:"createdAt"`
}

type ReferralAttribution struct {
	ID                          string                    `json:"id"`
	ProgramID                   string                    `json:"programId"`
	AccountID                   string                    `json:"accountId"`
	LinkID                      string                    `json:"linkId"`
	ReferrerCustomerID          string                    `json:"referrerCustomerId"`
	RefereeCustomerID           string                    `json:"refereeCustomerId"`
	Status                      ReferralAttributionStatus `json:"status"`
	ReferrerRewardCodeID        *string                   `json:"referrerRewardCodeId"`
	RefereeRewardCodeID         *string                   `json:"refereeRewardCodeId"`
	QualifyingCheckoutSessionID *string                   `json:"qualifyingCheckoutSessionId"`
	ExpiresAt                   string                    `json:"expiresAt"`
	ExternalSource              *string                   `json:"externalSource"`
	ExternalRef                 *string                   `json:"externalRef"`
	ClickedAt                   string                    `json:"clickedAt"`
	SignedUpAt                  *string                   `json:"signedUpAt"`
	RewardedAt                  *string                   `json:"rewardedAt"`
	VoidedAt                    *string                   `json:"voidedAt"`
	VoidReason                  *string                   `json:"voidReason"`
	CreatedAt                   string                    `json:"createdAt"`
}

type ProgramStats struct {
	TotalLinks        int     `json:"totalLinks"`
	TotalClicks       int     `json:"totalClicks"`
	TotalSignups      int     `json:"totalSignups"`
	TotalRewards      int     `json:"totalRewards"`
	AttributedRevenue float64 `json:"attributedRevenue"`
	ConversionRate    float64 `json:"conversionRate"`
}

type MyReward struct {
	Role         string       `json:"role"` // "referrer" | "referee"
	AttributionID string      `json:"attributionId"`
	Code         string       `json:"code"`
	DiscountType DiscountType `json:"discountType"`
	Value        float64      `json:"value"`
	Currency     string       `json:"currency"`
	ExpiresAt    *string      `json:"expiresAt"`
	Redeemed     bool         `json:"redeemed"`
	Active       bool         `json:"active"`
	EarnedAt     string       `json:"earnedAt"`
}

type IssueLinkInput struct {
	AccountID  string `json:"accountId"`
	CustomerID string `json:"customerId"`
}

type IssueLinkResult struct {
	Link *ReferralLink `json:"link"`
}

type ResolveLinkResult struct {
	Link ReferralLink `json:"link"`
}

type RecordClickInput struct {
	AccountID string `json:"accountId"`
	Code      string `json:"code"`
}

type RecordClickResult struct {
	LinkID    string `json:"linkId"`
	ProgramID string `json:"programId"`
}

type AttributeSignupInput struct {
	AccountID         string  `json:"accountId"`
	RefereeCustomerID string  `json:"refereeCustomerId"`
	RefereeEmail      string  `json:"refereeEmail"`
	ReferrerEmail     *string `json:"referrerEmail,omitempty"`
	LinkCode          string  `json:"linkCode"`
	ExternalSource    *string `json:"externalSource,omitempty"`
	ExternalRef       *string `json:"externalRef,omitempty"`
}

type AttributeSignupResult struct {
	Attribution *ReferralAttribution `json:"attribution"`
}

type AttributeCheckoutStartInput struct {
	AccountID         string `json:"accountId"`
	CustomerID        string `json:"customerId"`
	CheckoutSessionID string `json:"checkoutSessionId"`
}

type AttributeCheckoutStartResult struct {
	Stamped       bool   `json:"stamped"`
	AttributionID string `json:"attributionId,omitempty"`
}

type FulfillRewardInput struct {
	CheckoutSessionID string  `json:"checkoutSessionId"`
	Status            string  `json:"status"`
	Currency          string  `json:"currency"`
	Amount            float64 `json:"amount"`
}

type FulfillRewardResult struct {
	Issued             bool   `json:"issued"`
	Reason             string `json:"reason,omitempty"`
	ReferrerCodeID     string `json:"referrerCodeId,omitempty"`
	RefereeCodeID      string `json:"refereeCodeId,omitempty"`
}

type VoidAttributionResult struct {
	Voided     bool `json:"voided"`
	ClawedBack bool `json:"clawedBack"`
}

type MyRewardsListWrap struct {
	Items []MyReward `json:"items"`
}

type ExpirePendingResult struct {
	Expired int `json:"expired"`
}

// ─── Partner billing ───────────────────────────────────────────────────

type PartnerWorkspace struct {
	AccountID      string  `json:"accountId"`
	Partner        string  `json:"partner"`
	DiscountRate   float64 `json:"discountRate"`
	BrandName      *string `json:"brandName"`
	BusinessEmail  *string `json:"businessEmail"`
	CreatedAt      string  `json:"createdAt"`
}

type ProvisionWorkspaceInput struct {
	AccountID     string  `json:"accountId"`
	Partner       string  `json:"partner"` // "storlaunch"
	DiscountRate  float64 `json:"discountRate"`
	BrandName     *string `json:"brandName,omitempty"`
	BusinessEmail *string `json:"businessEmail,omitempty"`
}

type PartnerUsageSummary struct {
	Partner string `json:"partner"`
	Period  struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"period"`
	Totals struct {
		Redemptions      int `json:"redemptions"`
		ReferralRewards  int `json:"referralRewards"`
		Reminders        int `json:"reminders"`
		ChargeableCents  int `json:"chargeableCents"`
	} `json:"totals"`
	ByMerchant []struct {
		AccountID       string `json:"accountId"`
		Redemptions     int    `json:"redemptions"`
		ReferralRewards int    `json:"referralRewards"`
		Reminders       int    `json:"reminders"`
		ChargeableCents int    `json:"chargeableCents"`
	} `json:"byMerchant"`
}

type PartnerUsageQuery struct {
	Partner string `query:"partner"`
	From    string `query:"from"`
	To      string `query:"to"`
}

// ─── Webhook event envelope ────────────────────────────────────────────

// WebhookEvent is the parsed envelope returned by VerifyWebhook. Data
// is left as raw JSON so callers can decode into a concrete type.
type WebhookEvent struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	OccurredAt string          `json:"occurredAt"`
	AccountID  *string         `json:"accountId"`
	Data       json.RawMessage `json:"data"`
	Metadata   map[string]any  `json:"metadata"`
}

// ─── Misc list/wrap helpers ────────────────────────────────────────────

type DeletedResult struct {
	Deleted bool `json:"deleted"`
}

type RevokedResult struct {
	Revoked bool `json:"revoked"`
}

type ApplicableListWrap struct {
	Items []PublicApplicableCode `json:"items"`
}
