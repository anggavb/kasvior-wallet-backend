package controller

import (
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kasvior-wallet-backend/internal/helper"
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
	claims, ok := helper.CheckClaims(ctx)
	if !ok {
		return
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		helper.JSONBadRequest(ctx)
		return
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		helper.JSONBadRequest(ctx)
		return
	}

	search := strings.TrimSpace(ctx.DefaultQuery("search", ""))

	res, err := tc.transactionService.FindReceivers(ctx.Request.Context(), claims.UserId, search, page, limit)
	if err != nil {
		log.Println("Error: ", err.Error())
		helper.JSONInternalServerError(ctx)
		return
	}

	helper.JSONSuccess(ctx, res, "Get Receivers Successfully")
}
