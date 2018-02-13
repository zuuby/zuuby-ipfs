package server

import (
	"io"
	"log"
	"net/http"
	"os"

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

	http.HandleFunc("/", res.HandleIndex)

	return &res
}

func (s *Server) Serve() {

	listenString := s.ListenHost + ":" + s.ListenPort
	s.Logger.Println("Serving http://" + listenString)
	s.Logger.Fatal(http.ListenAndServe(listenString, nil))
}

func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {

  s.Logger.Printf("%v", *r)
	s.httpHeaders(w)
	io.WriteString(w, "hello, world<br/><br/>")
}

func (s *Server) httpHeaders(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
}
