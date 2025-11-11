package query

import (
	"context"

	"github.com/MaisamV/wallet/internal/probes/entity"
	"github.com/MaisamV/wallet/platform/logger"
)

// Static ping response to avoid creating new instances on every request
var staticPingResponse = &entity.PingResponse{
	Message: "PONG",
}

// PingQueryHandler handles ping queries
type PingQueryHandler struct {
	logger logger.Logger
}

// NewPingQueryHandler creates a new ping query handler
func NewPingQueryHandler(logger logger.Logger) *PingQueryHandler {
	return &PingQueryHandler{
		logger: logger,
	}
}

// Handle processes the ping query and returns a ping response
func (h *PingQueryHandler) Handle(ctx context.Context) (*entity.PingResponse, error) {
	h.logger.Debug().Msg("Processing ping request")
	// Return the static response for better performance
	return staticPingResponse, nil
}
