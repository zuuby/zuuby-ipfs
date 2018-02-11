package request

type CRUDVerb int

const GET CRUDVerb = 0
const PUT CRUDVerb = 1

type Request struct {
  Verb CRUDVerb
  Payload []byte
}

type RequestChan chan *Request
