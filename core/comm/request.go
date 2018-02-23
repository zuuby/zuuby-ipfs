package comm

type verb int

const Get verb = 0
const Put verb = 1
const Delete verb = 2

type Push struct {
	Verb    verb
	Payload []byte
}

func NewPush(v verb, pl []byte) *Push {
	return &Push{v, pl}
}

type Request struct {
	*Push
	Res chan *Response
}

func NewGetReq(payload []byte) *Request {
	return &Request{
		NewPush(Get, payload),
		make(chan *Response, 1),
	}
}

func NewPutReq(payload []byte) *Request {
	return &Request{
		NewPush(Put, payload),
		make(chan *Response, 1),
	}
}

func NewDeleteReq(payload []byte) *Request {
	return &Request{
		NewPush(Delete, payload),
		make(chan *Response, 1),
	}
}
