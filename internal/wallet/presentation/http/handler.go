package http

import (
	"github.com/MaisamV/wallet/internal/wallet/application/command"
	"github.com/MaisamV/wallet/internal/wallet/application/query"
	"github.com/MaisamV/wallet/internal/wallet/presentation/dto"
	"github.com/MaisamV/wallet/platform/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid/v5"
	"net/http"
	"strconv"
)

type WalletHandler struct {
	logger                 logger.Logger
	withdrawHandler        *command.WithdrawCommandHandler
	chargeHandler          *command.ChargeCommandHandler
	balanceHandler         *query.GetBalanceQueryHandler
	transactionPageHandler *query.GetTransactionPageQueryHandler
}

func NewWalletHandler(logger logger.Logger, withdrawHandler *command.WithdrawCommandHandler,
	chargeHandler *command.ChargeCommandHandler, balanceHandler *query.GetBalanceQueryHandler,
	transactionPageHandler *query.GetTransactionPageQueryHandler) *WalletHandler {
	return &WalletHandler{
		logger:                 logger,
		withdrawHandler:        withdrawHandler,
		chargeHandler:          chargeHandler,
		balanceHandler:         balanceHandler,
		transactionPageHandler: transactionPageHandler,
	}
}

// RegisterRoutes registers the documentation routes
func (h *WalletHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/api/v1/wallet")
	h.logger.Info().Msg("Registering wallet routes")
	group.Get("/:userid", h.GetBalance)
	group.Get("/:userid/transactions", h.GetTransactions)
	group.Post("/:userid/withdraw", h.Withdraw)
	group.Post("/:userid/charge", h.Charge)
	h.logger.Info().Msg("wallet routes registered successfully")
}

func (h *WalletHandler) GetBalance(c *fiber.Ctx) error {
	ctx := c.Context()

	userID, err := strconv.ParseInt(c.Params("userid"), 10, 64)
	if err != nil {
		return h.respondError(c, http.StatusBadRequest, err, "Could not parse userid")
	}

	q := query.GetBalanceQuery{UserID: userID}
	balance, err := h.balanceHandler.Handle(ctx, q)
	if err != nil {
		return h.respondError(c, http.StatusInternalServerError, err, "Could not fetch user's wallet balance")
	}

	return c.Status(http.StatusOK).JSON(dto.ToResponse(balance))
}

func (h *WalletHandler) GetTransactions(c *fiber.Ctx) error {
	ctx := c.Context()

	userID, err := strconv.ParseInt(c.Params("userid"), 10, 64)
	if err != nil {
		return h.respondError(c, http.StatusBadRequest, err, "Could not parse userid")
	}
	limit, err := strconv.ParseInt(c.Params("limit", "10"), 10, 64)
	if err != nil {
		return h.respondError(c, http.StatusBadRequest, err, "Could not parse limit")
	}

	var cursor *uuid.UUID
	cursorString := c.Query("cursor", "")
	if cursorString != "" {
		crs, err := uuid.FromString(cursorString)
		if err != nil {
			return h.respondError(c, http.StatusBadRequest, err, "Could not parse cursor")
		}
		cursor = &crs
	}

	q := query.GetTransactionPageQuery{
		UserID: userID,
		Cursor: cursor,
		Limit:  int(limit),
	}
	transactionPage, err := h.transactionPageHandler.Handle(ctx, q)
	if err != nil {
		return h.respondError(c, http.StatusInternalServerError, err, "Could not fetch transactions")
	}

	return c.Status(http.StatusOK).JSON(dto.ToResponse(transactionPage))
}

func (h *WalletHandler) Withdraw(c *fiber.Ctx) error {
	ctx := c.Context()

	userID, err := strconv.ParseInt(c.Params("userid"), 10, 64)
	if err != nil {
		return h.respondError(c, http.StatusBadRequest, err, "Could not parse userid")
	}
	withdraw := dto.Transaction{}
	err = c.BodyParser(&withdraw)
	if err != nil {
		return h.respondError(c, http.StatusBadRequest, err, "Could not parse the json")
	}

	idempotency, err := uuid.FromString(withdraw.Idempotency)
	if err != nil {
		return h.respondError(c, http.StatusBadRequest, err, "Could not parse idempotency")
	}

	cmd := command.WithdrawCommand{
		UserId:      userID,
		Amount:      withdraw.Amount,
		Idempotency: &idempotency,
		ReleaseTime: withdraw.ReleaseTime,
	}
	transactionID, err := h.withdrawHandler.Handle(ctx, cmd)
	if err != nil {
		return h.respondError(c, http.StatusInternalServerError, err, "Could not withdraw")
	}

	return c.Status(http.StatusOK).JSON(dto.ToResponse(transactionID.String()))
}

func (h *WalletHandler) Charge(c *fiber.Ctx) error {
	ctx := c.Context()

	userID, err := strconv.ParseInt(c.Params("userid"), 10, 64)
	if err != nil {
		return h.respondError(c, http.StatusBadRequest, err, "Could not parse userid")
	}
	charge := dto.Transaction{}
	err = c.BodyParser(&charge)
	if err != nil {
		return h.respondError(c, http.StatusBadRequest, err, "Could not parse the json")
	}

	idempotency, err := uuid.FromString(charge.Idempotency)
	if err != nil {
		return h.respondError(c, http.StatusBadRequest, err, "Could not parse idempotency")
	}

	cmd := command.ChargeCommand{
		UserId:      userID,
		Amount:      charge.Amount,
		Idempotency: &idempotency,
		ReleaseTime: charge.ReleaseTime,
	}
	transactionID, err := h.chargeHandler.Handle(ctx, cmd)
	if err != nil {
		return h.respondError(c, http.StatusInternalServerError, err, "Could not charge")
	}

	return c.Status(http.StatusOK).JSON(dto.ToResponse(transactionID.String()))
}

func (h *WalletHandler) respondError(c *fiber.Ctx, status int, err error, message string) error {
	h.logger.Error().Err(err).Msg(message)
	return c.Status(status).JSON(dto.ToErrorWithMessage(err, message))
}
