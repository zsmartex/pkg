package pkg

import (
	"github.com/shopspring/decimal"

	"github.com/zsmartex/pkg/order"
)

type PayloadAction = string

var (
	ActionSubmit        PayloadAction = "submit"
	ActionCancel        PayloadAction = "cancel"
	ActionCancelWithKey PayloadAction = "cancel_with_key" // this method will not notify to the user
	ActionReload        PayloadAction = "reload"
	ActionNew           PayloadAction = "new"
)

type MatchingPayloadMessage struct {
	Action PayloadAction   `json:"action"`
	Order  *order.Order    `json:"order"`
	Key    *order.OrderKey `json:"key"`
}

type GetDepthPayload struct {
	Market string `json:"market"`
	Limit  int64  `json:"limit"`
}

type DepthJSON struct {
	Asks     [][]decimal.Decimal `json:"asks"`
	Bids     [][]decimal.Decimal `json:"bids"`
	Sequence int64               `json:"sequence"`
}

type EnqueueEventKind string

var (
	EnqueueEventKindPublic  EnqueueEventKind = "public"
	EnqueueEventKindPrivate EnqueueEventKind = "private"
)
