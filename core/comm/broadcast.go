package comm

type SignalChan chan struct{}
type ReadOnlySignalChan <-chan struct{}

type BroadcastChan chan interface{}
type ReadOnlyBroadcastChan <-chan interface{}

func NewSignal() SignalChan {
  return make(SignalChan, 0)
}

// TODO: implement general broadcast channel
// func NewBroadcast() SignalChan {}
