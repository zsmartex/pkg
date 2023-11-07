package twilio_fx

import (
	"testing"

	"github.com/zsmartex/pkg/v2/config"
)

func TestSendMessage(t *testing.T) {
	client := New(twilioParams{
		Config: config.Twilio{
			PhoneNumber: "+16106248045",
			AccountSid:  "ACb8e046899dc2b3de11f42e1da1371692",
			AuthToken:   "ea00a3e882732b1d7a190fe10ee95abe",
		},
	})

	if err := client.SendSMS("+84835553240", "hello"); err != nil {
		t.Error(err)
	}
}
