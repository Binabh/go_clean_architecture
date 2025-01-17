package setup

import (
	"clean-architecture/domain"
	"clean-architecture/pkg"
	"clean-architecture/pkg/infrastructure"
	"clean-architecture/pkg/middlewares"
	"clean-architecture/pkg/services"

	"go.uber.org/fx"
)

var TestModule = fx.Options(
	services.Module,
	infrastructure.Module,
	middlewares.Module,
	pkg.Module,
	domain.Module,
	fx.Decorate(TestEnvReplacer),
)
