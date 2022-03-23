package services

import (
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/vault/api"
)

type TOTPService struct {
	vault                  *api.Client
	application_name       string
	vault_application_name string
}

func NewTOTPService(vault_addr, token, application_name, vault_application_name string) (*TOTPService, error) {
	config := &api.Config{
		Address: vault_addr,
		Timeout: time.Second * 2,
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetToken(token)

	ts := &TOTPService{
		vault:                  client,
		application_name:       application_name,
		vault_application_name: vault_application_name,
	}

	ts.startRenewToken(token)

	return ts, nil
}

func (t *TOTPService) startRenewToken(token string) error {
	secret, err := t.vault.Auth().Token().Lookup(token)
	if err != nil {
		return err
	}

	var renewable bool
	if v, ok := secret.Data["renewable"]; ok {
		renewable, _ = v.(bool)
	}

	if !renewable {
		return nil
	}

	watcher, err := t.vault.NewLifetimeWatcher(&api.LifetimeWatcherInput{
		Secret: &api.Secret{
			Auth: &api.SecretAuth{
				Renewable:   renewable,
				ClientToken: token,
			},
		},
	})

	if err != nil {
		return err
	}

	go watcher.Start()
	go func() {
		for {
			select {
			case <-watcher.DoneCh():
				return
			case <-watcher.RenewCh():
			}
		}
	}()
	return nil
}

func (t *TOTPService) Create(uid, email string) {
	t.vault.Logical().Write(t.totp_key(uid), map[string]interface{}{
		"generate":     true,
		"issuer":       os.Getenv("BARONG_APP_NAME"),
		"account_name": email,
		"qr_size":      100,
	})
}

func (t *TOTPService) Validate(uid, code string) bool {
	secret, err := t.vault.Logical().Write(t.totp_code_key(uid), map[string]interface{}{
		"code": code,
	})
	if err != nil {
		return false
	}

	return secret.Data["valid"].(bool)
}

func (t *TOTPService) Delete(uid string) {
	t.vault.Logical().Delete(t.totp_key(uid))
}

func (t *TOTPService) Exist(uid string) bool {
	secret, err := t.vault.Logical().Read(t.totp_key(uid))
	if err != nil {
		return false
	}

	return secret != nil
}

func (t *TOTPService) get_vault_application_name() string {
	return os.Getenv("BARONG_VAULT_APPLICATION_NAME")
}

func (t *TOTPService) totp_key(uid string) string {
	return fmt.Sprintf("totp/keys/%s_%s", t.get_vault_application_name(), uid)
}

func (t *TOTPService) totp_code_key(uid string) string {
	return fmt.Sprintf("totp/codes/%s_%s", t.get_vault_application_name(), uid)
}
