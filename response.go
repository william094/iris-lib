package iris_lib

import "github.com/kataras/iris/v12"

type ErrMsg struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"-"`
}

type Result struct {
	MsgCode      int         `json:"msg_code"`
	ErrorCode    int         `json:"error_code"`
	ErrorMessage string      `json:"error_message"`
	Data         interface{} `json:"data"`
}

func OK(data interface{}) *Result {
	return &Result{
		MsgCode:      0,
		ErrorCode:    0,
		ErrorMessage: "Success",
		Data:         data,
	}
}

func Error(err ErrMsg) *Result {
	return &Result{
		MsgCode:      err.Code,
		ErrorCode:    err.Code,
		ErrorMessage: err.Msg,
		Data:         nil,
	}
}

func Return(result *Result, ctx iris.Context) {
	ctx.Record()
	ctx.JSON(result)
}
