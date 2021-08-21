package pkg

import (
	"github.com/shopspring/decimal"

	"github.com/zsmartex/pkg/order"
)

type PayloadAction = string

var (
	ActionSubmit         PayloadAction = "submit"
	ActionCancel         PayloadAction = "cancel"
	ActionCancelWithKey  PayloadAction = "cancel_with_key" // this method will not notify to the user
	ActionReload         PayloadAction = "reload"
	ActionNew            PayloadAction = "new"
	ActionGetMarketPrice PayloadAction = "get_market_price"
)

type MatchingPayloadMessage struct {
	Action PayloadAction   `json:"action"`
	Order  *order.Order    `json:"order"`
	Key    *order.OrderKey `json:"key"`
	Market string          `json:"market"`
}

type GetDepthPayload struct {
	Market string `json:"market"`
	Limit  int    `json:"limit"`
}

type DepthJSON struct {
	Asks     [][]decimal.Decimal `json:"asks"`
	Bids     [][]decimal.Decimal `json:"bids"`
	Sequence uint64              `json:"sequence"`
}
