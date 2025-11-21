package response

//
//import (
//	"net/http"
//
//	"github.com/gin-gonic/gin"
//)
//
//type Body struct {
//	Code    int    `json:"code"`
//	Msg     string `json:"msg"`
//	Data    any    `json:"data"`
//	TraceID string `json:"traceID"`
//}
//
//func Result(c *gin.Context, code int, data any, err any) {
//	if data == nil {
//		data = make(map[string]string, 0)
//	}
//	resp := Body{
//		Code: code,
//		Data: data,
//		Msg:  Msg[code],
//	}
//
//	if e, ok := err.(error); ok {
//		resp.Msg = resp.Msg + ": " + e.Error()
//	}
//	c.JSON(http.StatusOK, resp)
//}
//
//func Success(c *gin.Context, data any) {
//	Result(c, Succ, data, nil)
//}
//
//func Fail(c *gin.Context, code int, err any) {
//	Result(c, code, nil, err)
//}
