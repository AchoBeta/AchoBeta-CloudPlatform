package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SystemJsonMsgResponse struct {
	Ctx *gin.Context
}

type JsonMsgResult struct {
	Code    int
	Message string
	Data    interface{}
}

const SUCCESS_CODE = 200
const SUCCESS_MSG = "成功"
const ERROR_MSG = "错误"

func NewResponse(c *gin.Context) *SystemJsonMsgResponse {
	return &SystemJsonMsgResponse{Ctx: c}
}

func (r *SystemJsonMsgResponse) Success(data interface{}) {
	res := JsonMsgResult{}
	res.Code = SUCCESS_CODE
	res.Message = SUCCESS_MSG
	res.Data = data
	r.Ctx.JSON(http.StatusOK, res)
}

func (r *SystemJsonMsgResponse) Error(mc MsgCode) {
	r.error(mc.Code, mc.Msg)
}

func (r *SystemJsonMsgResponse) error(code int, message string) {
	if message == "" {
		message = ERROR_MSG
	}
	res := JsonMsgResult{}
	res.Code = code
	res.Message = message
	res.Data = gin.H{}
	r.Ctx.JSON(http.StatusOK, res)
}
