package controller

import (
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

func (tc *TransactionController) FindReceivers(ctx *gin.Context) {
	claims, ok := jwttoken.CheckClaims(ctx)
	if !ok {
		return
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		response.JSONBadRequest(ctx)
		return
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		response.JSONBadRequest(ctx)
		return
	}

	search := strings.TrimSpace(ctx.DefaultQuery("search", ""))

	res, err := tc.transactionService.FindReceivers(ctx.Request.Context(), claims.UserId, search, page, limit)
	if err != nil {
		log.Println("Error: ", err.Error())
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONSuccess(ctx, res, "Get Receivers Successfully")
}

func (tc *TransactionController) CreateTopup(ctx *gin.Context) {
	claims, ok := jwttoken.CheckClaims(ctx)
	if !ok {
		return
	}

	var body dto.TopupRequest
	if !binder.BindFormat(ctx, &body, binding.JSON) {
		return
	}

	paymentMethod, err := tc.transactionService.CreateTransactionWithDetails(ctx.Request.Context(), claims.UserId, "topup", body)
	if err != nil {
		log.Println("Error: ", err.Error())
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONCreated(ctx, dto.TopupResponse{
		Amount:        body.Amount,
		PaymentMethod: paymentMethod,
		Discount:      body.Discount,
		Tax:           body.Tax,
		SubTotal:      body.SubTotal,
	}, "Topup Successfully!")
}
