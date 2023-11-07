package vault_fx

import "go.uber.org/fx"

var (
	Module = fx.Module("vault_fx", vaultProviders)

	vaultProviders = fx.Provide(
		New,
		NewTransit,
		NewTOTP,
	)
)
