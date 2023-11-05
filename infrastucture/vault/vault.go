package vault

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/vault/api"
)

type VaultService struct {
	client *api.Client
}

func New(vaultAddr, token string) (*VaultService, error) {
	config := &api.Config{
		Address: vaultAddr,
		Timeout: time.Second * 2,
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetToken(token)

	vs := &VaultService{
		client: client,
	}

	vs.startRenewToken(token)

	return vs, nil
}

func (s *VaultService) startRenewToken(token string) error {
	secret, err := s.client.Auth().Token().Lookup(token)
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

	watcher, err := s.client.NewLifetimeWatcher(&api.LifetimeWatcherInput{
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
	return s.client.Logical().Read(path)
}

func (s *VaultService) Write(path string, data map[string]interface{}) (*api.Secret, error) {
	return s.client.Logical().Write(path, data)
}

func (s *VaultService) Unwrap(path string) (*api.Secret, error) {
	return s.client.Logical().Unwrap(path)
}

func (s *VaultService) Delete(path string) (*api.Secret, error) {
	return s.client.Logical().Delete(path)
}

func (s *VaultService) Health(ctx context.Context) error {
	res, err := s.client.Sys().Health()
	if err != nil {
		return err
	}

	if res.Sealed {
		return errors.New("vault is sealed")
	}

	if !res.Initialized {
		return errors.New("vault is not initialized")
	}

	return nil
}
