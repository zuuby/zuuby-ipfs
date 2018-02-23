package client

import (
	"fmt"
	"os/exec"

	"github.com/zuuby/zuuby-ipfs/core/comm"
)

type WorkerPool struct {
	Size int
	stop comm.ReadOnlySignalChan
	done chan bool
}

func New(stopChan comm.ReadOnlySignalChan, n int) WorkerPool {
	return WorkerPool{
		Size: n,
		stop: stopChan,
		done: make(chan bool, n),
	}
}

func (w WorkerPool) Start() chan *comm.Request {
	fmt.Printf("[worker] Starting a pool of %d workers\n", w.Size)

	// TODO: I don't actually know what to make the RequestChan capacity
	rqch := make(chan *comm.Request, w.Size*w.Size)

	for i := 0; i < w.Size; i++ { // spawn n workers
		go func() {
			fmt.Println("[worker] Worker started.")
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

func process(r *comm.Request) {
	var res *comm.Response
	switch r.Verb {
	case comm.Get:
		res = get(r.Payload)
	case comm.Put:
		res = put(r.Payload)
	default:
		fmt.Printf("[worker] Unknown request verb: %v\n", r.Verb)
		res = comm.NewBadRequest()
	}

	r.Res <- res
}

func get(hash []byte) *comm.Response {
	cmd := exec.Command("sh", "-c", "ipfs cat "+string(hash))
	out, err := cmd.Output()

	if err != nil {
		fmt.Printf("[worker] Get(%s) request stopped with error: %v", string(hash), err)
		return comm.NewServerError()
	}

	fmt.Printf("[worker] %s", string(out))
	return comm.NewSuccess(out)
}

func put(contents []byte) *comm.Response {
	strcmd := fmt.Sprintf("echo \"%s\" | ipfs add -Q", string(contents))
	out, err := exec.Command("sh", "-c", strcmd).Output()

	if err != nil {
		fmt.Printf("[worker] Put request stopped with error: %v", err)
		return comm.NewServerError()
	}

	// TODO: must schedule hash for pinning
	fmt.Printf("[worker] %s", string(out))
	return comm.NewSuccess(out)
}
