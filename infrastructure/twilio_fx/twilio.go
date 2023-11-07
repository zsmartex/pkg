package twilio_fx

import (
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/zsmartex/pkg/v2/config"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"twilio_fx",
		fx.Provide(
			New,
		),
	)
)

type Client struct {
	client      *twilio.RestClient
	phoneNumber string
}

type twilioParams struct {
	fx.In

	Config config.Twilio
}

func New(params twilioParams) *Client {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: params.Config.AccountSid,
		Password: params.Config.AuthToken,
	})

	return &Client{
		phoneNumber: params.Config.PhoneNumber,
		client:      client,
	}
}

func (t *Client) SendSMS(to_number string, text string) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(to_number)
	params.SetFrom(t.phoneNumber)
	params.SetBody(text)

	if _, err := t.client.Api.CreateMessage(params); err != nil {
		return err
	}

	return nil
}
