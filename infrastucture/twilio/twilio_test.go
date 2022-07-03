package twilio

import "testing"

func TestSendMessage(t *testing.T) {
	client := New("+16106248045", "ACb8e046899dc2b3de11f42e1da1371692", "ea00a3e882732b1d7a190fe10ee95abe")

	if err := client.SendSMS("+84835553240", "hello"); err != nil {
		t.Error(err)
	}
}
