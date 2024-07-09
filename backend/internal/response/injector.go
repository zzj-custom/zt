package response

type Reply struct {
	Code   ErrCode `json:"code"`
	Msg    string  `json:"msg"`
	Result any     `json:"result,omitempty"`
}

func (r *Reply) GetCode() ErrCode {
	return r.Code
}

func (r *Reply) GetMsg() string {
	return r.Msg
}

type IReply interface {
	GetCode() int
	GetMsg() string
}

func OkReply(r any) *Reply {
	return &Reply{
		Code:   Ok,
		Msg:    Ok.String(),
		Result: r,
	}
}

func FailReply(code ErrCode) *Reply {
	return &Reply{
		Code:   code,
		Msg:    code.String(),
		Result: nil,
	}
}

func FailReplyWithResult(code ErrCode, msg string, r any) *Reply {
	return &Reply{
		Code:   code,
		Msg:    msg,
		Result: r,
	}
}
