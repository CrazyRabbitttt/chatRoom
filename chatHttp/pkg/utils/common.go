package utils

import "github.com/gin-gonic/gin"

// common files in utils package

const (
	Ok      = 0
	Invalid = 1
	Online  = 2
	OffLine = 3
)

type CommonResponse struct {
	StatusCode int32
	StatusMsg  string
}

// 从中间件中获得之前存储的东西
func GetId(ctx *gin.Context) int64 {
	val, _ := ctx.Get("user_id")
	userId, ok := val.(int64)
	if !ok {
		return -1
	}
	return userId
}
