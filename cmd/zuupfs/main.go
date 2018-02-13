package main

import (
  "fmt"
  "os"
  "os/signal"

  "github.com/zuuby/zuuby-ipfs/daemon"
  "github.com/zuuby/zuuby-ipfs/client"
  "github.com/zuuby/zuuby-ipfs/core/comm"
)

func main() {
  fmt.Println("Starting the zuupfs daemon")

  // create the Daemon and start it
  dmn := daemon.New()
  stop, _ := dmn.Start() // returns a stop channel

  dmn.WaitReady() // Have to wait for the daemon to start before we send requests

  // create a worker pool and start it
  wp := client.New(stop, 5)
  rc := wp.Start()

  svr := server.NewHttpServer(httpPort, rc)
	svr.Serve()

  res := make(chan string)

  // rc <- &request.Request{
  //   Verb: request.PUT,
  //   Payload: []byte("Some sting data"),
  //   Response: res,
  // }

  // After you run the program once, copy the output of the add command and
  // paste it in the payload below to test and confirm. Uncomment this block
  // and comment the above put request.
  // TODO: do this in code so we don't have to manually
  rc <- &comm.Request{
    Verb: comm.GET,
    Payload: []byte("QmXyTGQm8p7QnQWWMo3yeowTVSEypFj7ibty441BkfowZs"),
    Response: res,
  }

  defer func() {
    svr.Stop()    // close the api endpoints
    dmn.Stop()    // will signal the workers
    wp.WaitDone() // wait for workers to stop
  }()

  c := make(chan os.Signal, 1)   // signal channel
  signal.Notify(c, os.Interrupt) // get interrupt(^C) signal
  <-c                            // block until interrupt
}
