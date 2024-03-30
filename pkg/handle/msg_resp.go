package handle

import (
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JsonMsgResponse struct {
	Ctx *app.RequestContext
}

type JsonMsgResult struct {
	Code    int
	Message string
	Data    interface{}
}

const SUCCESS_CODE = 200
const SUCCESS_MSG = "成功"
const ERROR_MSG = "错误"

func NewResponse(c *app.RequestContext) *JsonMsgResponse {
	return &JsonMsgResponse{Ctx: c}
}

func (r *JsonMsgResponse) Success(data interface{}) {
	res := JsonMsgResult{}
	res.Code = SUCCESS_CODE
	res.Message = SUCCESS_MSG
	res.Data = data
	r.Ctx.JSON(http.StatusOK, res)
}

func (r *JsonMsgResponse) Error(mc MsgCode) {
	r.error(mc.Code, mc.Msg)
}

func (r *JsonMsgResponse) error(code int, message string) {
	if message == "" {
		message = ERROR_MSG
	}
	res := JsonMsgResult{}
	res.Code = code
	res.Message = message
	res.Data = gin.H{}
	r.Ctx.JSON(http.StatusOK, res)
}
