package client

import (
    "fmt"
    "os/exec"

    "github.com/zuuby/zuuby-ipfs/core/broadcast"
    "github.com/zuuby/zuuby-ipfs/core/request"
)

type RequestPtr *request.Request
type RequestChan request.RequestChan
type ReceiverSignalChan broadcast.ReceiverSignalChan
type WorkerDoneChan chan bool

type WorkerPool struct {
  Size int                // num workers
  stop ReceiverSignalChan // stop signal from daemon
  done WorkerDoneChan     // channel for workers to signal their shutdown
}

func New(stopChan ReceiverSignalChan, n int) WorkerPool {
  return WorkerPool{
    Size: n,
    stop: stopChan,
    done: make(chan bool, n),
  }
}

func (w WorkerPool) Start() RequestChan {
  fmt.Printf("[worker] Starting a pool of %d workers\n", w.Size)

  // TODO: I don't actually know what to make the RequestChan capacity
  rqch := make(RequestChan, w.Size * w.Size)

  for i := 0; i < w.Size; i++ { // spawn n workers
    go func () {
      for {
        select {
        case r := <-rqch: // process a request
          go process(r)
        case <-w.stop: // the daemon has stopped
          fmt.Println("[worker] Daemon stopped. Stopping worker")
          w.done <- true
          return
        }
      }
    }()
  }

  return rqch
}

func (w WorkerPool) WaitDone() {
  for i := 0; i < w.Size; i++ {
    <-w.done
  }
}

func process(r RequestPtr) {
  switch r.Verb {
  case request.GET:
    get(r.Payload)
  case request.PUT:
    put(r.Payload)
  default:
    fmt.Printf("[worker] Unknown request verb: %v\n", r.Verb)
  }
}

func get(hash []byte) {
  cmd := exec.Command("ipfs", "object", "get", string(hash))
  out, err := cmd.Output()

  if err != nil {
    fmt.Printf("[worker] Get(%s) request stopped with error: %v", string(hash), err)
    return
  }

  // TODO: must return the returned object
  fmt.Printf("[worker] %s", string(out))
}

func put(contents []byte) {
  strcmd := fmt.Sprintf("echo \"%s\" | ipfs add -Q", string(contents))
  //fmt.Println(strcmd)
  out, err := exec.Command("sh", "-c", strcmd).Output()

  if err != nil {
    fmt.Printf("[worker] Put request stopped with error: %v", err)
    return
  }

  // TODO: must return hash somehow and schedule it for pinning
  fmt.Printf("[worker] %s", string(out))
}
