package comm

const SuccessCode int = 200

// TODO: Add more descriptive errors
const ClientError int = 400
const ServerError int = 500

const SuccessMsg string = "Success"
const BadRequestMsg string = "Bad Request"
const NotFoundMsg string = "Not Found"
const ServerMsg string = "Internal Server Error"

type Response struct {
	Msg     string `json:"message"`
	Code    int    `json:"code"`
	Payload []byte `json:"data"`
}

func NewSuccess(payload []byte) *Response {
	return &Response{
		Msg:     SuccessMsg,
		Code:    SuccessCode,
		Payload: payload,
	}
}

func NewBadRequest() *Response {
	return &Response{
		Msg:     BadRequestMsg,
		Code:    ClientError,
		Payload: make([]byte, 0),
	}
}

func NewNotFound() *Response {
	return &Response{
		Msg:     NotFoundMsg,
		Code:    ClientError,
		Payload: make([]byte, 0),
	}
}

func NewServerError() *Response {
	return &Response{
		Msg:     ServerMsg,
		Code:    ServerError,
		Payload: make([]byte, 0),
	}
}
