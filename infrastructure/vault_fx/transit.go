package vault_fx

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/zsmartex/pkg/v2/config"
	"go.uber.org/fx"
)

type Transit struct {
	client          *Client
	applicationName string
}

type TransitParams struct {
	fx.In

	Client *Client
	Config config.Vault
}

func NewTransit(params TransitParams) *Transit {
	return &Transit{
		client:          params.Client,
		applicationName: params.Config.ApplicationName,
	}
}

func (s *Transit) Encrypt(key, value string) (*api.Secret, error) {
	valueString := base64.StdEncoding.EncodeToString([]byte(value))

	return s.client.Write(fmt.Sprintf("transit/encrypt/%s_%s", s.applicationName, key), map[string]interface{}{
		"plaintext": valueString,
	})
}

func (s *Transit) Decrypt(key, ciphertext string) (*api.Secret, error) {
	return s.client.Write(fmt.Sprintf("transit/decrypt/%s_%s", s.applicationName, key), map[string]interface{}{
		"ciphertext": ciphertext,
	})
}

func (s *Transit) Delete(key string) (*api.Secret, error) {
	return s.client.Delete(fmt.Sprintf("transit/encrypt/%s_%s", s.applicationName, key))
}
