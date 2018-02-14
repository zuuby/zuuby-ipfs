package server

import (
	"io"
	"log"
	"net/http"
	"os"
  "encoding/json"

  "github.com/zuuby/zuuby-ipfs/core/comm"
)

type Server struct {
	ListenHost string
	ListenPort string
	Logger     *log.Logger
  requestChan comm.RequestChan
}

func NewHttpServer(httpPort string, rc comm.RequestChan) *Server {

	res := Server{
		ListenHost: "localhost",
		ListenPort: httpPort,
		Logger:     log.New(os.Stdout, "server> ", log.Ltime|log.Lshortfile),
    requestChan: rc,
  }

  http.HandleFunc("/api/get", res.HandleGet)
  http.HandleFunc("/api/put", res.HandlePut)

	return &res
}

func (s *Server) Serve() {

	listenString := s.ListenHost + ":" + s.ListenPort
	s.Logger.Println("Serving http://" + listenString)
	s.Logger.Fatal(http.ListenAndServe(listenString, nil))
}

func (s *Server) HandleGet(w http.ResponseWriter, r *http.Request) {
  s.Logger.Printf("%v", *r)
  if r.Method != "GET" || r.Method != "" {
    BadRequestError(w)
    return
  }

  req := NewGetReq([]byte(""))
  s.requestChan <- req

  res := <-req.Res

	s.httpHeaders(w)
	Success(w, res)
}

func (s *Server) HandlePut(w http.ResponseWriter, r *http.Request) {
  s.Logger.Printf("%v", *r)
  if r.Method != "PUT" {
    BadRequestError(w)
    return
  }

  req := NewPutReq([]byte("Some string data"))
  s.requestChan <- req

  res := req.Res

	s.httpHeaders(w)
	Success(w, res)
}

func (s *Server) httpHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
}

func Success(w http.ResponseWriter, res *comm.Response) {
  data, _ = json.Marshal(res)
  io.WriteString(w, data)
}

func BadRequestError(w http.ResponseWriter) {
  data, _ = json.Marshal(comm.NewBadRequest())
  io.WriteString(w, data)
}

func NotFoundError(w http.ResponseWriter) {
  data, _ = json.Marshal(comm.NewNotFound())
  io.WriteString(w, data)
}

func ServerError(w http.ResponseWriter) {
  data, _ = json.Marshal(comm.NewServerError())
  io.WriteString(w, data)
}
