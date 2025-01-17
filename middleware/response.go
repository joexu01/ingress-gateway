package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joexu01/ingress-gateway/lib"
	"strings"
)

type ResponseCode int

//1000以下为通用码，1000以上为用户自定义码
const (
	SuccessCode ResponseCode = iota
	UndefErrorCode
	ValidErrorCode
	InternalErrorCode

	InvalidRequestErrorCode ResponseCode = 401
	CustomizeCode           ResponseCode = 1000

	GroupAllSaveFlowError ResponseCode = 2001
)

type Response struct {
	ErrorCode ResponseCode `json:"errno"`
	ErrorMsg  string       `json:"errmsg"`
	Data      interface{}  `json:"data"`
	TraceId   interface{}  `json:"trace_id"`
	Stack     interface{}  `json:"stack"`
}

func ResponseError(c *gin.Context, code ResponseCode, err error) {
	trace, _ := c.Get("trace")
	traceContext, _ := trace.(*lib.TraceContext)
	traceId := ""
	if traceContext != nil {
		traceId = traceContext.TraceId
	}

	stack := ""
	if c.Query("is_debug") == "1" || lib.GetConfEnv() == "dev" {
		stack = strings.Replace(fmt.Sprintf("%+v", err), err.Error()+"\n", "", -1)
	}

	resp := &Response{ErrorCode: code, ErrorMsg: err.Error(), Data: "", TraceId: traceId, Stack: stack}

	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
	c.AbortWithStatusJSON(500, resp)
	//c.JSON(500, resp)
	//if code <= 1000 {
	//	c.JSON(500, resp)
	//} else {
	//	c.JSON(500, resp)
	//}
	//if code <= 1000 {
	//	_ = c.AbortWithError(500, err)
	//} else {
	//	_ = c.AbortWithError(500, err)
	//}
	//TODO: No receiver here at first
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	trace, _ := c.Get("trace")
	traceContext, _ := trace.(*lib.TraceContext)
	traceId := ""
	if traceContext != nil {
		traceId = traceContext.TraceId
	}

	resp := &Response{ErrorCode: SuccessCode, ErrorMsg: "", Data: data, TraceId: traceId}
	c.JSON(200, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
}
