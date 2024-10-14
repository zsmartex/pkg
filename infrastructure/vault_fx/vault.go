package vault_fx

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/zsmartex/pkg/v2/config"
	"go.uber.org/fx"
)

type Client struct {
	client *api.Client
}

type vaultParams struct {
	fx.In

	Config config.Vault
}

func New(params vaultParams) (*Client, error) {
	config := &api.Config{
		Address: params.Config.Address,
		Timeout: time.Second * 2,
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetToken(params.Config.Token)

	cli := &Client{
		client: client,
	}

	err = cli.startRenewToken(params.Config.Token)
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func (s *Client) startRenewToken(token string) error {
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

func (s *Client) Read(path string) (*api.Secret, error) {
	return s.client.Logical().Read(path)
}

func (s *Client) Write(path string, data map[string]interface{}) (*api.Secret, error) {
	return s.client.Logical().Write(path, data)
}

func (s *Client) Unwrap(path string) (*api.Secret, error) {
	return s.client.Logical().Unwrap(path)
}

func (s *Client) Delete(path string) (*api.Secret, error) {
	return s.client.Logical().Delete(path)
}

func (s *Client) Health(ctx context.Context) error {
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

func (s *Client) HAStatus() (*api.HAStatusResponse, error) {
	return s.client.Sys().HAStatus()
}
