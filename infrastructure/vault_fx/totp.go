package vault_fx

import (
	"fmt"

	"go.uber.org/fx"

	"github.com/zsmartex/pkg/v2/config"
)

type TOTP struct {
	client          *Client
	applicationName string
}

type totpParams struct {
	fx.In

	Client *Client
	Config config.Vault
}

func NewTOTP(params totpParams) *TOTP {
	return &TOTP{
		client:          params.Client,
		applicationName: params.Config.ApplicationName,
	}
}

func (s *TOTP) Create(issuer, uid, email string) (map[string]interface{}, error) {
	if result, err := s.client.Write(s.totpKey(uid), map[string]interface{}{
		"generate":     true,
		"issuer":       s.applicationName,
		"account_name": email,
		"qr_size":      100,
	}); err != nil {
		return nil, err
	} else {
		return result.Data, nil
	}
}

func (s *TOTP) SetApplicationName(name string) {
	s.applicationName = name
}

func (s *TOTP) Validate(uid, code string) bool {
	secret, err := s.client.Write(s.totpCodeKey(uid), map[string]interface{}{
		"code": code,
	})
	if err != nil {
		return false
	}

	return secret.Data["valid"].(bool)
}

func (s *TOTP) Delete(uid string) {
	s.client.Delete(s.totpKey(uid))
}

func (s *TOTP) Exist(uid string) bool {
	secret, err := s.client.Read(s.totpKey(uid))
	if err != nil {
		return false
	}

	return secret != nil
}

func (s *TOTP) totpKey(uid string) string {
	return fmt.Sprintf("totp/keys/%s_%s", s.applicationName, uid)
}

func (s *TOTP) totpCodeKey(uid string) string {
	return fmt.Sprintf("totp/code/%s_%s", s.applicationName, uid)
}
