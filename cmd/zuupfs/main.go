package main

import (
  "fmt"
  "os"
  "os/signal"

  "github.com/zuuby/zuuby-ipfs/daemon"
  "github.com/zuuby/zuuby-ipfs/client"
  "github.com/zuuby/zuuby-ipfs/core/request"
)

//type Daemon daemon.Daemon
//type WorkerPool client.WorkerPool

func main() {
  fmt.Println("Starting the zuupfs daemon")

  // create the Daemon and start it
  dmn := daemon.New()
  stop, _ := dmn.Start() // returns a stop channel
  dmn.WaitReady()

  // create a worker pool and start it
  wp := client.New(stop, 5)
  rc := wp.Start()

  // rc <- &request.Request{
  //   Verb: request.PUT,
  //   Payload: []byte("Some sting data"),
  // }

  // After you run the program once, copy the output of the add command and
  // paste it in the payload below to test and confirm. Uncomment this block
  // and comment the above put request.
  // TODO: do this in code so we don't have to manually
  rc <- &request.Request{
    Verb: request.GET,
    Payload: []byte("QmXyTGQm8p7QnQWWMo3yeowTVSEypFj7ibty441BkfowZs"),
  }

  defer func() {
    dmn.Stop()    // will signal the workers
    wp.WaitDone() // wait for workers to stop
  }()

  c := make(chan os.Signal, 1)   // signal channel
  signal.Notify(c, os.Interrupt) // get interrupt(^C) signal
  <-c                            // block until interrupt
}
