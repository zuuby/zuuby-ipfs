package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/zuuby/zuuby-ipfs/core/comm"
)

const timeout_after time.Duration = 10

type Server struct {
	ListenHost  string
	ListenPort  string
	Logger      *log.Logger // TODO remove this
	requestChan chan *comm.Request
}

func NewHttpServer(httpPort string, rc chan *comm.Request) *Server {

	svr := Server{
		ListenHost:  "localhost",
		ListenPort:  httpPort,
		Logger:      log.New(os.Stdout, "server> ", log.Ltime|log.Lshortfile),
		requestChan: rc,
	}

	http.HandleFunc("/", svr.handleIndex)
	http.HandleFunc("/api/file", svr.handleApi)

	return &svr
}

func (s *Server) Serve() {
	listenString := ":" + s.ListenPort
	s.Logger.Println("[server] Serving http://" + listenString)
	go func() {
		s.Logger.Fatal(http.ListenAndServe(listenString, nil))
	}()
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	badRequestError(w)
}

func (s *Server) handleApi(w http.ResponseWriter, r *http.Request) {

	s.Logger.Printf("%v", *r)
	setHttpHeaders(w)

	if r.Method == "GET" {
		s.handleGet(w, r)
	} else if r.Method == "PUT" || r.Method == "POST" {
		s.handlePut(w, r)
	} else if r.Method == "DELETE" {

		fmt.Println("[server] DELETE Method not handled yet.")
		badRequestError(w)
		//TODO: s.handleDelete(w, r)
	} else {

		fmt.Printf("[server] Unknown method %s.", r.Method)
		badRequestError(w)
	}
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		badRequestError(w)
		return
	}

	hash := r.URL.Query().Get("hash")

	if hash == "" {
		s.Logger.Printf("[server] Hash parameter missing.")
		badRequestError(w)
		return
	}

	req := comm.NewGetReq([]byte(hash))
	s.requestChan <- req

	select {
	case res := <-req.Res: // wait for response
		handleResponse(w, res)
	case <-time.After(timeout_after * time.Second):
		fmt.Println("[server] Client timeout. Failing request")
		serverError(w)
	}
}

func (s *Server) handlePut(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		badRequestError(w)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		serverError(w)
		return
	}

	req := comm.NewPutReq([]byte(body))
	s.requestChan <- req

	select {
	case res := <-req.Res: // wait for response
		handleResponse(w, res)
	case <-time.After(timeout_after * time.Second):
		fmt.Println("[server] Client timeout. Failing request")
		serverError(w)
	}
}

func handleResponse(w http.ResponseWriter, res *comm.Response) {
	w.WriteHeader(int(res.Code))
	data, _ := json.Marshal(res)
	fmt.Fprintln(w, string(data))
}

func setHttpHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func badRequestError(w http.ResponseWriter) {
	http.Error(w, comm.BadRequestMsg, int(comm.ClientError))
}

func notFoundError(w http.ResponseWriter) {
	http.Error(w, comm.NotFoundMsg, int(comm.ClientError))
}

func serverError(w http.ResponseWriter) {
	http.Error(w, comm.ServerMsg, int(comm.ServerError))
}
