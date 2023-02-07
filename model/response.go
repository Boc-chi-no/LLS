package model

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response This structure represents the standard response payload
type Response struct {
	Code  int `json:"code"`
	Data interface{} `json:"data"`
	Detail string `json:"detail"`
	Fail bool `json:"fail"`
	Message string `json:"message"`
	Success bool `json:"success"`
	Type string `json:"type"`
}

func SuccessResponse(c *gin.Context,data map[string]interface{}){
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Data: data,
		Detail: "",
		Fail: false,
		Message: "",
		Success: true,
		Type: "",
	})
}

func FailureResponse(c *gin.Context,httpCode int,sysCode int, message string,detail string) {
	c.AbortWithStatusJSON(httpCode, Response{
		Code: sysCode,
		Data: "",
		Detail: detail,
		Fail: true,
		Message: message,
		Success: false,
		Type: "",
	})
}
