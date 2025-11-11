package swagger

import (
	swaggerQuery "github.com/MaisamV/wallet/internal/swagger/application/query"
	"github.com/MaisamV/wallet/internal/swagger/infrastructure"
	swaggerHttp "github.com/MaisamV/wallet/internal/swagger/presentation/http"
	"github.com/MaisamV/wallet/platform/config"
	"github.com/MaisamV/wallet/platform/logger"
	"github.com/google/wire"
)

// ProvideSwaggerLoader provides a swagger loader
func ProvideSwaggerLoader(logger logger.Logger, config *config.Config) (*infrastructure.SwaggerLoader, error) {
	loader := infrastructure.NewSwaggerLoader(logger, config.Swagger)
	if err := loader.Init(); err != nil {
		return nil, err
	}
	return loader, nil
}

// ProvideSwaggerQueryHandler provides a swagger query handler
func ProvideSwaggerQueryHandler(logger logger.Logger, swaggerLoader *infrastructure.SwaggerLoader) *swaggerQuery.SwaggerQueryHandler {
	return swaggerQuery.NewSwaggerQueryHandler(logger, swaggerLoader)
}

// ProvideDocsHandler provides a docs HTTP handler
func ProvideDocsHandler(logger logger.Logger, swaggerQueryHandler *swaggerQuery.SwaggerQueryHandler) *swaggerHttp.DocsHandler {
	return swaggerHttp.NewDocsHandler(logger, swaggerQueryHandler)
}

// SwaggerSet is a wire provider set for all swagger dependencies
var SwaggerSet = wire.NewSet(
	ProvideSwaggerLoader,
	ProvideSwaggerQueryHandler,
	ProvideDocsHandler,
)
