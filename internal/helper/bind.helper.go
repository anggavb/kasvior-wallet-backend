package helper

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func BindFormat(ctx *gin.Context, requestData any, binder binding.Binding) bool {
	if err := ctx.ShouldBindWith(&requestData, binder); err != nil {
		log.Println("Error", err.Error())
		JSONBadRequest(ctx)
		return false
	}

	return true
}
