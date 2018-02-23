package comm

type SignalChan chan struct{}

func (sc SignalChan) Signal() {
	close(sc)
}

type ReadOnlySignalChan <-chan struct{}

// TODO: implement general broadcast channel
// type BroadcastChan chan interface{}
// type ReadOnlyBroadcastChan <-chan interface{}

func NewSignalChan() SignalChan {
	return make(SignalChan, 0)
}

// TODO: implement general broadcast channel
// func NewBroadcast() SignalChan {}
