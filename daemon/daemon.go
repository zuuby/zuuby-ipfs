package daemon

import (
    "errors"
    "fmt"
    "os/exec"
    "syscall"
    "time"

    "github.com/zuuby/zuuby-ipfs/core/broadcast"
)

type SenderSignalChan broadcast.SenderSignalChan

type Daemon struct {
  cmd *exec.Cmd
}

func New() Daemon {
  return Daemon{
    cmd: exec.Command("/usr/local/bin/ipfs", "daemon"),
  }
}

// NOTE: ReceiverSignalChan is a read-only chan and SenderSignalChan is read/write
// HACK: compiler won't allow return type written as broadcast.ReceiverSignalChan
func (d Daemon) Start() (<-chan struct{}, error) {
  stopChan := make(SenderSignalChan, 0)

  fmt.Println("[daemon] Starting the ipfs daemon")

	err := d.cmd.Start()
  if err != nil {
    fmt.Printf("[daemon] Daemon failed to start with error: %v\n", err)
    return nil, err
  }

  // started successfully
  go func() {
    if err := d.cmd.Wait(); err != nil {
      fmt.Printf("[daemon] The ipfs daemon exitted with error %v.\n", err)
    }
    fmt.Println("[daemon] The ipfs daemon has stopped. Broadcasting stop signal.")
    close(stopChan)
  }()

  return stopChan, nil
}

func (d Daemon) Stop() error {
  fmt.Println("[daemon] Stopping the daemon")
  if d.cmd.Process != nil {

    // cmd.ProcessState is non-nil when cmd.Wait() has been called
    // If that is the case, we should check that it hasn't already stopped
    if d.cmd.ProcessState != nil && d.cmd.ProcessState.Exited() {
      fmt.Println("[daemon] Daemon process already exited.")
      return nil
    }

    fmt.Println("[daemon] Sending SIGINT")
    if err := d.cmd.Process.Signal(syscall.SIGINT); err != nil {
      fmt.Printf("[daemon] Error sending signal %v.\n", err)
      return err
    }
  }
  return nil
}

// HACK: There is currently no good way to determine if the daemon is running AND
// ready to take commands
func (d Daemon) WaitReady() error {
  timeouts := 0
  for {
    timeouts += 1
    if timeouts > 10 {
      fmt.Println("[daemon] Daemon timeout.")
      return errors.New("daemon: Too many timeouts. Daemon not ready.")
    }
    cmd := exec.Command("/usr/local/bin/ipfs", "cat", "/ipfs/QmYwAPJzv5CZsnA625s3Xf2nemtYgPpHdWEz79ojWnPbdG/readme")
    _, err := cmd.Output()
    if err != nil {
      fmt.Println("[daemon] Daemon not ready yet ...")
      time.Sleep(1 * time.Second)
      continue
    }

    fmt.Println("[daemon] Daemon ready!")
    return nil
  }
}
