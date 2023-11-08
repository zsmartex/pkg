package rango_fx

import (
	"go.uber.org/fx"
)

var Module = fx.Module("rango", fx.Provide(New))
