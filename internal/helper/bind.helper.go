package helper

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// using generics data type for body parameter (belajar dari kelasnya pak eko pzn)
func BindFormat[T any](ctx *gin.Context, requestData *T, binder binding.Binding) bool {
	if err := ctx.ShouldBindWith(&requestData, binder); err != nil {
		log.Println("Error", err.Error())
		JSONBadRequest(ctx)
		return false
	}

	return true
}
