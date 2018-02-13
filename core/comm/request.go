package comm

type CRUDVerb int

const GET CRUDVerb = 0
const PUT CRUDVerb = 1
const DELETE CRUDVerb = 2

type Push struct {
  Verb CRUDVerb
  Payload []byte
}

type PushChan chan *Push

func NewPush(v CRUDVerb, pl []byte) *Push {
  return &Push{v, pl}
}

type Request struct {
  *Push
  Res ResponseChan
}

type RequestChan chan *Request

func NewGetReq(payload []byte) *Request {
  return &Request{
    NewPush(GET, payload),
    make(ResponseChan, 1),
  }
}

func NewPutReq(payload []byte) *Request {
  return &Request{
    NewPush(PUT, payload),
    make(ResponseChan, 1),
  }
}
