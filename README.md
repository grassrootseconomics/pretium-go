# pretium-go

![GitHub Tag](https://img.shields.io/github/v/tag/grassrootseconomics/pretium-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/grassrootseconomics/pretium-go.svg)](https://pkg.go.dev/github.com/grassrootseconomics/pretium-go)

Go wrapper for the Pretium API.

## Installation

```bash
$ go get github.com/grassrootseconomics/pretium-go
```

## Usage

```go
client := pretium.New(settlementAddress, apiKey, callbackURL)

resp, err := client.Onramp(ctx, pretium.KES, pretium.OnrampBody{
    Shortcode:     "0712345678",
    Amount:        100,
    MobileNetwork: pretium.SAFARICOM,
    Chain:         pretium.CELO,
    Asset:         pretium.CUSD,
    Address:       "0x...",
})
```

### Handling webhooks

Pretium emits three webhook variants — an offramp payout notification, an
onramp payment confirmation, and an onramp asset-release notification — all
posted to the `callback_url` configured on the client. `WebhookPayload`
exposes the union of all fields; use `Event()` to discriminate.

```go
http.HandleFunc("/pretium/webhook", func(w http.ResponseWriter, r *http.Request) {
    payload, err := pretium.ParseWebhookRequest(r)
    if err != nil {
        http.Error(w, "bad payload", http.StatusBadRequest)
        return
    }

    switch payload.Event() {
    case pretium.WebhookEventAssetReleased:
        // Onramp asset release: payload.IsReleased and payload.TransactionHash are set.
    case pretium.WebhookEventStatus:
        // Offramp payout or onramp payment confirmation.
        if payload.Status == pretium.StatusComplete {
            // settle on your side
        }
    }

    w.WriteHeader(http.StatusOK)
})
```

## License

[AGPL-3.0](LICENSE)
