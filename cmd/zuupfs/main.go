package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/zuuby/zuuby-ipfs/client"
	"github.com/zuuby/zuuby-ipfs/daemon"
	//"github.com/zuuby/zuuby-ipfs/core/comm"
	"github.com/zuuby/zuuby-ipfs/core/server"
)

const port = "5000"

func main() {
	fmt.Println("Starting the zuupfs daemon")

	// create the Daemon and start it
	dmn := daemon.New()
	stop, _ := dmn.Start()

	dmn.WaitReady() // Have to wait for the daemon to start before we send requests

	// create a worker pool and start it
	wp := client.New(stop, 5)
	rc := wp.Start()

	// create a server
	svr := server.NewHttpServer(port, rc)
	svr.Serve()

	defer func() {
		//svr.Stop()    // TODO close the api endpoints
		dmn.Stop()    // will signal the workers
		wp.WaitDone() // wait for workers to stop
	}()

	c := make(chan os.Signal, 1)   // signal channel
	signal.Notify(c, os.Interrupt) // get interrupt(^C) signal
	<-c                            // block until interrupt
}
