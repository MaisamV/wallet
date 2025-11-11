package infrastructure

import (
	"fmt"
	"github.com/MaisamV/wallet/platform/config"
	"os"

	"github.com/MaisamV/wallet/platform/logger"
)

// SwaggerLoader implements the ports.SwaggerProvider interface
type SwaggerLoader struct {
	logger          logger.Logger
	openApiFilePath string
	swaggerFilePath string
	openapi         []byte
	swaggerHtml     []byte
}

// NewSwaggerLoader creates a new instance of DocsAdapter
func NewSwaggerLoader(logger logger.Logger, cfg config.SwaggerConfig) *SwaggerLoader {
	return &SwaggerLoader{
		logger:          logger,
		openApiFilePath: cfg.OpenApiFilePath,
		swaggerFilePath: cfg.SwaggerFilePath,
	}
}

func (a *SwaggerLoader) Init() error {
	a.logger.Info().Str("openapi_path", a.openApiFilePath).Msg("Loading OpenAPI specification")
	data, err := os.ReadFile(a.openApiFilePath)
	if err != nil {
		a.logger.Error().Err(err).Str("path", a.openApiFilePath).Msg("Failed to read OpenAPI spec")
		return fmt.Errorf("failed to read OpenAPI spec: %w", err)
	}
	a.openapi = data
	a.logger.Info().Int("size_bytes", len(data)).Msg("OpenAPI specification loaded successfully")

	a.logger.Info().Str("swagger_path", a.swaggerFilePath).Msg("Loading Swagger HTML")
	html, err := os.ReadFile(a.swaggerFilePath)
	if err != nil {
		a.logger.Error().Err(err).Str("path", a.swaggerFilePath).Msg("Failed to read Swagger HTML")
		return fmt.Errorf("failed to read Swagger HTML: %w", err)
	}
	a.swaggerHtml = html
	a.logger.Info().Int("size_bytes", len(html)).Msg("Swagger HTML loaded successfully")
	return nil
}

// GetOpenAPISpec loads the OpenAPI specification from file
func (a *SwaggerLoader) GetOpenAPISpec() ([]byte, error) {
	a.logger.Debug().Msg("Serving OpenAPI specification")
	return a.openapi, nil
}

// GetSwaggerHTML generates the Swagger UI HTML
func (a *SwaggerLoader) GetSwaggerHTML() ([]byte, error) {
	a.logger.Debug().Msg("Serving Swagger UI HTML")
	return a.swaggerHtml, nil
}
