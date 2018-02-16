package comm

type StatusCode int

const SUCCESS_CODE StatusCode = 200

// TODO: Add more descriptive errors
const CLIENT_ERROR StatusCode = 400
const SERVER_ERROR StatusCode = 500

type StatusMsg string

const SUCCESS_MSG = "Success"
const BAD_REQUEST_MSG = "Bad Request"
const NOT_FOUND_MSG = "Not Found"
const SERVER_MSG = "Internal Server Error"

type Response struct {
  Msg StatusMsg   `json:"message"`
  Code StatusCode `json:"code"`
  Payload []byte  `json:"data"`
}

type ResponseChan chan *Response

func NewSuccess(payload []byte) *Response {
  return &Response {
    Msg:      SUCCESS_MSG,
    Code:     SUCCESS_CODE,
    Payload:  payload,
  }
}

func NewBadRequest() *Response {
  return &Response{
    Msg:      BAD_REQUEST_MSG,
    Code:     CLIENT_ERROR,
    Payload:  make([]byte, 0),
  }
}

func NewNotFound() *Response {
  return &Response{
    Msg:      NOT_FOUND_MSG,
    Code:     CLIENT_ERROR,
    Payload:  make([]byte, 0),
  }
}

func NewServerError() *Response {
  return &Response{
    Msg:      SERVER_MSG,
    Code:     SERVER_ERROR,
    Payload:  make([]byte, 0),
  }
}
