# ripllo (Go)

Official Go SDK for [Ripllo](https://ripllo.com) — discounts,
referrals, abandoned cart, pixels, feeds, blog, creator marketplace,
marketing tools (channels, contacts, contact-lists, campaigns,
funnels, inbox, segments), insights, partner admin.

Mirrors `@forjio/ripllo-node` 0.2.2 and `ripllo` (PyPI) 0.1.x at the
API-surface level.

## Install

```bash
go get github.com/hachimi-cat/ripllo-go
```

## Auth

Pattern 2 HMAC partner-billing — same scheme storlaunch + fulkruma
use. Set credentials via env or pass them to `NewClient`:

```bash
RIPLLO_KEY_ID=AKIARPLO...
RIPLLO_SECRET=...
RIPLLO_BASE_URL=https://ripllo.com   # optional override
```

```go
import (
    "context"
    ripllo "github.com/hachimi-cat/ripllo-go"
)

func main() {
    rip, err := ripllo.NewClient(ripllo.ClientOptions{})
    if err != nil { panic(err) }

    // Platform-admin keys can act on behalf of a merchant:
    scoped := rip.ForMerchant("acc_storlaunch_merchant_123")

    page, err := scoped.DiscountCodes.List(context.Background(), ripllo.DiscountCodesListParams{
        Limit: ripllo.IntPtr(20),
    })
    if err != nil { /* *ripllo.Error */ }
    _ = page
}
```

## Webhook verification

```go
import (
    "io"
    "net/http"
    "os"

    ripllo "github.com/hachimi-cat/ripllo-go"
)

func handler(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    ev, err := ripllo.VerifyWebhook(
        body,
        r.Header.Get("Ripllo-Signature"),
        os.Getenv("RIPLLO_WEBHOOK_SECRET"),
        nil,
    )
    if err != nil {
        http.Error(w, "bad signature", 400)
        return
    }
    _ = ev // ev.Type, ev.Data (json.RawMessage), ev.AccountID
    w.WriteHeader(204)
}
```

## Notes

- Path-without-querystring signing: the SDK signs the path portion
  only (matches Node 0.2.x + Python 0.1.x bug fix; mismatch on `?limit=` etc.
  was the original incident).
- Stdlib only — no third-party dependencies.
- `Passthrough` for partners that need to relay arbitrary
  merchant-portal requests without hand-typed methods.
- Untyped marketing routes return the alias `JSON` (= `map[string]any`)
  — same shape as the Node SDK's `Record<string, unknown>` returns.

## License

MIT
