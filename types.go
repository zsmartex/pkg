package pkg

import "github.com/zsmartex/pkg/order"

type PayloadAction = string

var (
	ActionSubmit        PayloadAction = "submit"
	ActionCancel        PayloadAction = "cancel"
	ActionCancelWithKey PayloadAction = "cancel_with_key" // this method will not notify to the user
	ActionReload        PayloadAction = "reload"
)

type MatchingPayloadMessage struct {
	Action PayloadAction   `json:"action"`
	Order  *order.Order    `json:"order"`
	Key    *order.OrderKey `json:"key"`
	Market string          `json:"market"`
}
