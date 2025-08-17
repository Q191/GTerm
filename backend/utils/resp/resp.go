package resp

import (
	"github.com/Q191/GTerm/backend/consts/messages"
)

type Resp struct {
	Ok   bool   `json:"ok"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
	Code string `json:"code"`
}

func result(ok bool, code string, data any) *Resp {
	resp := &Resp{
		Ok:   ok,
		Code: code,
		Data: data,
	}

	if code != "" && messages.CodeMapping != nil {
		if defaultMsg, exists := messages.CodeMapping[code]; exists {
			resp.Msg = defaultMsg
		}
	}

	return resp
}

func Ok() *Resp {
	return result(true, "", nil)
}

func OkWithCode(code string) *Resp {
	return result(true, code, nil)
}

func OkWithData(data any) *Resp {
	return result(true, "", data)
}

func OkWithCodeAndData(code string, data any) *Resp {
	return result(true, code, data)
}

func Fail() *Resp {
	return result(false, "", nil)
}

func FailWithCode(code string) *Resp {
	return result(false, code, nil)
}

func FailWithData(data any) *Resp {
	return result(false, "", data)
}

func FailWithCodeAndData(code string, data any) *Resp {
	return result(false, code, data)
}

func OkWithMsg(msg string) *Resp {
	return result(true, msg, nil)
}

func FailWithMsg(msg string) *Resp {
	return result(false, msg, nil)
}

func OkWithDetailed(msg string, data any) *Resp {
	return result(true, msg, data)
}

func FailWithDetailed(msg string, data any) *Resp {
	return result(false, msg, data)
}
