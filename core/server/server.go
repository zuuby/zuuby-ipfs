package server

import (
	"io"
  "io/ioutil"
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
	s.Logger.Println("[server] Serving http://" + listenString)
  go func () {
    s.Logger.Fatal(http.ListenAndServe(listenString, nil))
  }()
}

func (s *Server) HandleGet(w http.ResponseWriter, r *http.Request) {
  s.Logger.Printf("%v", *r)
  if r.Method != "GET" {
    BadRequestError(w)
    return
  }

  hash := r.URL.Query().Get("hash")

  if hash == "" {
    s.Logger.Printf("[server] Hash parameter missing.")
    BadRequestError(w)
    return
  }

  req := comm.NewGetReq([]byte(hash))
  s.requestChan <- req

  res := <-req.Res

	s.httpHeaders(w)
	Success(w, res)
}

func (s *Server) HandlePut(w http.ResponseWriter, r *http.Request) {
  s.Logger.Printf("%v", *r)
  if r.Method != "POST" {
    BadRequestError(w)
    return
  }

  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    ServerError(w)
    return
  }

  req := comm.NewPutReq([]byte(body))
  s.requestChan <- req

  res := <-req.Res

	s.httpHeaders(w)
	Success(w, res)
}

func (s *Server) httpHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
}

func Success(w http.ResponseWriter, res *comm.Response) {
  data, _ := json.Marshal(res)
  io.WriteString(w, string(data))
}

func BadRequestError(w http.ResponseWriter) {
  data, _ := json.Marshal(comm.NewBadRequest())
  io.WriteString(w, string(data))
}

func NotFoundError(w http.ResponseWriter) {
  data, _ := json.Marshal(comm.NewNotFound())
  io.WriteString(w, string(data))
}

func ServerError(w http.ResponseWriter) {
  data, _ := json.Marshal(comm.NewServerError())
  io.WriteString(w, string(data))
}
