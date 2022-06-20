package twilio

import (
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Twilio struct {
	phoneNumber string
	accountSID  string
	authToken   string
	serviceID   string
	client      *twilio.RestClient
}

func New(phone_number, account_sid, auth_token, service_id string) *Twilio {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: account_sid,
		Password: auth_token,
	})

	return &Twilio{
		phoneNumber: phone_number,
		accountSID:  account_sid,
		authToken:   auth_token,
		serviceID:   service_id,
		client:      client,
	}
}

func (t *Twilio) SendSMS(to_number string, text string) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(to_number)
	params.SetFrom(t.phoneNumber)
	params.SetBody(text)

	if _, err := t.client.Api.CreateMessage(params); err != nil {
		return err
	}

	return nil
}
