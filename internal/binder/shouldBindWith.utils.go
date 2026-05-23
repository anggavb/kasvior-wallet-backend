package binder

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func BindFormat(ctx *gin.Context, requestData any, binder binding.Binding) error {
	if err := ctx.ShouldBindWith(requestData, binder); err != nil {
		log.Println("Error", err.Error())
		return err
	}

	return nil
}
