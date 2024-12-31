package resp

type Resp struct {
	Ok   bool   `json:"ok"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func result(ok bool, msg string, data any) *Resp {
	return &Resp{
		ok,
		msg,
		data,
	}
}

func Ok() *Resp {
	return result(true, "", nil)
}

func OkWithMsg(msg string) *Resp {
	return result(true, msg, nil)
}

func OkWithData(data any) *Resp {
	return result(true, "", data)
}

func OkWithDetailed(msg string, data any) *Resp {
	return result(true, msg, data)
}

func Fail() *Resp {
	return result(false, "", nil)
}

func FailWithMsg(msg string) *Resp {
	return result(false, msg, nil)
}

func FailWithData(data any) *Resp {
	return result(false, "", data)
}

func FailWithDetailed(msg string, data any) *Resp {
	return result(false, msg, data)
}
