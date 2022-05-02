package services

import (
	"time"

	"github.com/hashicorp/vault/api"
)

type VaultService struct {
	vault *api.Client
}

func NewVaultService(vault_addr, token string) (*VaultService, error) {
	config := &api.Config{
		Address: vault_addr,
		Timeout: time.Second * 2,
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetToken(token)

	vs := &VaultService{
		vault: client,
	}

	vs.startRenewToken(token)

	return vs, nil
}

func (s *VaultService) startRenewToken(token string) error {
	secret, err := s.vault.Auth().Token().Lookup(token)
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

	watcher, err := s.vault.NewLifetimeWatcher(&api.LifetimeWatcherInput{
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

func (s *VaultService) Read(path string) (*api.Secret, error) {
	return s.vault.Logical().Read(path)
}

func (s *VaultService) Write(path string, data map[string]interface{}) (*api.Secret, error) {
	return s.vault.Logical().Write(path, data)
}

func (s *VaultService) Unwrap(path string) (*api.Secret, error) {
	return s.vault.Logical().Unwrap(path)
}

func (s *VaultService) Delete(path string) (*api.Secret, error) {
	return s.vault.Logical().Delete(path)
}
