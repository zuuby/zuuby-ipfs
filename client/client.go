package client

import (
    "fmt"
    "os/exec"

    "github.com/zuuby/zuuby-ipfs/core/comm"
)

type RequestPtr *comm.Request
type ResponsePtr *comm.Response
type ReadOnlySignalChan comm.ReadOnlySignalChan
type WorkerDoneChan chan bool

type WorkerPool struct {
  Size int                // num workers
  stop ReadOnlySignalChan // stop signal from daemon
  done WorkerDoneChan     // channel for workers to signal their shutdown
}

func New(stopChan ReadOnlySignalChan, n int) WorkerPool {
  return WorkerPool{
    Size: n,
    stop: stopChan,
    done: make(chan bool, n),
  }
}

func (w WorkerPool) Start() comm.RequestChan {
  fmt.Printf("[worker] Starting a pool of %d workers\n", w.Size)

  // TODO: I don't actually know what to make the RequestChan capacity
  rqch := make(comm.RequestChan, w.Size * w.Size)

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
  var res ResponsePtr
  switch r.Verb {
  case comm.GET:
    res = get(r.Payload)
  case comm.PUT:
    res = put(r.Payload)
  default:
    fmt.Printf("[worker] Unknown request verb: %v\n", r.Verb)
    res = comm.NewBadRequest()
  }

  r.Res <- res
}

func get(hash []byte) ResponsePtr {
  cmd := exec.Command("ipfs", "object", "get", string(hash))
  out, err := cmd.Output()

  if err != nil {
    fmt.Printf("[worker] Get(%s) request stopped with error: %v", string(hash), err)
    return comm.NewServerError()
  }

  // TODO: must return the returned object
  fmt.Printf("[worker] %s", string(out))
  return comm.NewSuccess(out)
}

func put(contents []byte) ResponsePtr {
  strcmd := fmt.Sprintf("echo \"%s\" | ipfs add -Q", string(contents))
  out, err := exec.Command("sh", "-c", strcmd).Output()

  if err != nil {
    fmt.Printf("[worker] Put request stopped with error: %v", err)
    return comm.NewServerError()
  }

  // TODO: must return hash somehow and schedule it for pinning
  fmt.Printf("[worker] %s", string(out))
  return comm.NewSuccess(out)
}
