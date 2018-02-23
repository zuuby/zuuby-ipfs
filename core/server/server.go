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

const TIMEOUT_AFTER time.Duration = 10

type Server struct {
	ListenHost  string
	ListenPort  string
	Logger      *log.Logger
	requestChan comm.RequestChan
}

func NewHttpServer(httpPort string, rc comm.RequestChan) *Server {

	res := Server{
		ListenHost:  "localhost",
		ListenPort:  httpPort,
		Logger:      log.New(os.Stdout, "server> ", log.Ltime|log.Lshortfile),
		requestChan: rc,
	}

	http.HandleFunc("/", res.HandleIndex)
	http.HandleFunc("/api/file", res.HandleApi)

	return &res
}

func (s *Server) Serve() {
	listenString := ":" + s.ListenPort
	s.Logger.Println("[server] Serving http://" + listenString)
	go func() {
		s.Logger.Fatal(http.ListenAndServe(listenString, nil))
	}()
}

func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	BadRequestError(w)
}

func (s *Server) HandleApi(w http.ResponseWriter, r *http.Request) {

	s.Logger.Printf("%v", *r)
	SetHttpHeaders(w)

	if r.Method == "GET" {
		s.HandleGet(w, r)
	} else if r.Method == "PUT" || r.Method == "POST" {
		s.HandlePut(w, r)
	} else if r.Method == "DELETE" {

		fmt.Println("[server] DELETE Method not handled yet.")
		BadRequestError(w)
		//TODO: HandleDelete(w, r)
	} else {

		fmt.Printf("[server] Unknown method %s.", r.Method)
		BadRequestError(w)
	}
}

func (s *Server) HandleGet(w http.ResponseWriter, r *http.Request) {

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

	select {
	case res := <-req.Res: // wait for response
		HandleResponse(w, res)
	case <-time.After(TIMEOUT_AFTER * time.Second):
		fmt.Println("[server] Client timeout. Failing request")
		ServerError(w)
	}
}

func (s *Server) HandlePut(w http.ResponseWriter, r *http.Request) {

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

	select {
	case res := <-req.Res: // wait for response
		HandleResponse(w, res)
	case <-time.After(TIMEOUT_AFTER * time.Second):
		fmt.Println("[server] Client timeout. Failing request")
		ServerError(w)
	}
}

func HandleResponse(w http.ResponseWriter, res *comm.Response) {
	w.WriteHeader(int(res.Code))
	data, _ := json.Marshal(res)
	fmt.Fprintln(w, string(data))
}

func SetHttpHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func BadRequestError(w http.ResponseWriter) {
	http.Error(w, comm.BAD_REQUEST_MSG, int(comm.CLIENT_ERROR))
}

func NotFoundError(w http.ResponseWriter) {
	http.Error(w, comm.NOT_FOUND_MSG, int(comm.CLIENT_ERROR))
}

func ServerError(w http.ResponseWriter) {
	http.Error(w, comm.SERVER_MSG, int(comm.SERVER_ERROR))
}
