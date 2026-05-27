package controller

import (
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kasvior-wallet-backend/internal/apperrors"
	"github.com/kasvior-wallet-backend/internal/binder"
	"github.com/kasvior-wallet-backend/internal/dto"
	"github.com/kasvior-wallet-backend/internal/jwttoken"
	"github.com/kasvior-wallet-backend/internal/response"
	"github.com/kasvior-wallet-backend/internal/service"
)

type TransactionController struct {
	transactionService *service.TransactionService
}

func NewTransactionController(transactionService *service.TransactionService) *TransactionController {
	return &TransactionController{
		transactionService: transactionService,
	}
}

// FindReceivers godoc
// @Summary		Find transfer receivers
// @Description	Get receiver suggestions for the authenticated user.
// @Tags			Transactions
// @Produce		json
// @Security		ApiKeyAuth
// @Param			search		query		string	false	"Receiver name or phone number search keyword"	maxlength(100)
// @Param			page		query		int		false	"Page number"	default(1)
// @Param			limit		query		int		false	"Items per page"	minimum(1)	maximum(100)	default(10)
// @Success		200			{object}	dto.Response	"Get Receivers Successfully"
// @Failure		400			{object}	dto.Response	"Bad request"
// @Failure		401			{object}	dto.Response	"Unauthorized"
// @Failure		500			{object}	dto.Response	"Internal server error"
// @Router			/transaction/transfer/receivers [get]
func (tc *TransactionController) FindReceivers(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		return
	}

	var query dto.FindReceiversQueryRequest
	if err := binder.BindFormat(ctx, &query, binding.Query); err != nil {
		response.JSONBadRequest(ctx)
		return
	}

	page := 1
	if query.Page != nil {
		page = *query.Page
	}

	limit := 10
	if query.Limit != nil {
		limit = *query.Limit
	}

	search := strings.TrimSpace(query.Search)

	res, err := tc.transactionService.FindReceivers(ctx.Request.Context(), claims.UserId, search, page, limit)
	if err != nil {
		log.Println("Error: ", err.Error())
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONSuccess(ctx, res, "Get Receivers Successfully")
}

// CreateTopup godoc
// @Summary		Create topup transaction
// @Description	Create a topup transaction for the authenticated user.
// @Tags			Transactions
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			request		body		dto.TopupRequest	true	"Topup request body"
// @Success		201			{object}	dto.Response		"Topup created"
// @Failure		400			{object}	dto.Response		"Bad request"
// @Failure		401			{object}	dto.Response		"Unauthorized"
// @Failure		422			{object}	dto.Response		"Validation error"
// @Failure		500			{object}	dto.Response		"Internal server error"
// @Router			/transaction/topup [post]
func (tc *TransactionController) CreateTopup(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		return
	}

	var body dto.TopupRequest
	if err := binder.BindFormat(ctx, &body, binding.JSON); err != nil {
		errorMessages := binder.FormatValidationError(err)
		if len(errorMessages) > 0 && errorMessages["error"] != "" {
			response.JSONBadRequest(ctx)
			return
		}
		response.JSONUnprocessableEntity(ctx, errorMessages)
		return
	}

	if err := tc.transactionService.CreateTransactionWithDetails(ctx.Request.Context(), claims.UserId, body); err != nil {
		if errors.Is(err, apperrors.InvalidSubtotal) {
			log.Println("Error: ", err.Error())
			response.JSONBadRequest(ctx)
			return
		}
		log.Println("Error: ", err.Error())
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONCreated(ctx, nil, "Topup successfully")
}

// CreateTransfer godoc
// @Summary		Create transfer transaction
// @Description	Create a pending transfer transaction for the authenticated user.
// @Tags			Transactions
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			request	body		dto.TransferRequest				true	"Transfer request body"
// @Success		201		{object}	dto.TransactionCreatedResponse	"Transfer pending"
// @Failure		400		{object}	dto.Response					"Bad request"
// @Failure		401		{object}	dto.Response					"Unauthorized"
// @Failure		422		{object}	dto.Response					"Validation error"
// @Failure		500		{object}	dto.Response					"Internal server error"
// @Router			/transaction/transfer [post]
func (tc *TransactionController) CreateTransfer(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		return
	}

	var body dto.TransferRequest
	if err := binder.BindFormat(ctx, &body, binding.JSON); err != nil {
		errorMessages := binder.FormatValidationError(err)
		if len(errorMessages) > 0 && errorMessages["error"] != "" {
			response.JSONBadRequest(ctx)
			return
		}
		response.JSONUnprocessableEntity(ctx, errorMessages)
		return
	}

	res, err := tc.transactionService.CreatePendingTransfer(ctx.Request.Context(), claims.UserId, body)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidRecipient) {
			response.JSONBadRequestWithMessage(ctx, "Invalid recipient wallet")
			return
		}
		if errors.Is(err, apperrors.ErrSelfTransfer) {
			response.JSONBadRequestWithMessage(ctx, "Cannot transfer to own wallet")
			return
		}

		log.Println("Error: ", err.Error())
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONCreated(ctx, res, "Transfer pending")
}
