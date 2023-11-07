package session_fx

import "go.uber.org/fx"

var (
	Module = fx.Module("session_fx", sessionProviders)

	sessionProviders = fx.Provide(NewStore)
)
