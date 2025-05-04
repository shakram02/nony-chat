package generic

import (
	"fmt"
	"sync/atomic"
)

type Transport[T any] interface {
	Read(p []T) (n int, err error)
	Write(p []T) (n int, err error)
}

type ProtocolLayer[THigher any, TLower any] struct {
	bufferSize int
	transport  Transport[TLower]

	tx <-chan THigher
	rx chan<- THigher

	receiveAdapter func([]TLower) THigher
	sendAdapter    func(THigher) []TLower

	done chan struct{}

	isRunning atomic.Bool
}

func New[THigher any, TLower any](
	transport Transport[TLower],
	bufferSize int,
	tx <-chan THigher,
	rx chan<- THigher,
	receiveAdapter func([]TLower) THigher,
	sendAdapter func(THigher) []TLower,
) *ProtocolLayer[THigher, TLower] {
	return &ProtocolLayer[THigher, TLower]{
		transport:      transport,
		bufferSize:     bufferSize,
		tx:             tx,
		rx:             rx,
		receiveAdapter: receiveAdapter,
		sendAdapter:    sendAdapter,
		done:           make(chan struct{}),
		isRunning:      atomic.Bool{},
	}
}

func (l *ProtocolLayer[THigher, TLower]) Start() {
	if l.isRunning.Load() {
		return
	}

	l.isRunning.Store(true)
	go l.start()
}

func (l *ProtocolLayer[THigher, TLower]) Stop() {
	if !l.isRunning.Load() {
		return
	}

	close(l.done)
	close(l.rx)

	l.isRunning.Store(false)
}

func (l *ProtocolLayer[THigher, TLower]) start() {
	lowerRx := make(chan THigher)

	go func() {
		buffer := make([]TLower, l.bufferSize)
		for {
			if l.isRunning.Load() {
				break
			}

			n, err := l.transport.Read(buffer)
			if err != nil {
				panic(fmt.Sprintf("Failed to read transport: %v", err))
			}
			// This might block because of nobody receiving from
			// higher rx channel.
			lowerRx <- l.receiveAdapter(buffer[:n])
		}
	}()

	for {
		select {
		case toBeSentPacket := <-l.tx:
			l.transport.Write(l.sendAdapter(toBeSentPacket))
		case lowerReceived := <-lowerRx:
			l.rx <- lowerReceived
		case <-l.done:
			l.isRunning.Store(false)
			return
		}
	}
}
