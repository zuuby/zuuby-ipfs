package broadcast

type SenderSignalChan chan struct{} // broadcast a signal to all receivers

type ReceiverSignalChan <-chan struct{} // read a broadcasted signal
