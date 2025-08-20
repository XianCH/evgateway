package gateway

import (
	"fmt"

	"github.com/x14n/evgateway/internal/protocol"
)

type HandlerFunc func(*Gateway, *Session, protocol.Frame) error

type Dispatcher struct {
	handlers map[byte]HandlerFunc
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[byte]HandlerFunc),
	}
}

func (d *Dispatcher) RegisterHandler(cmd byte, handler HandlerFunc) {
	d.handlers[cmd] = handler
}

func (d *Dispatcher) Dispatch(gw *Gateway, session *Session, frame protocol.Frame) {
	if handler, ok := d.handlers[frame.Cmd]; ok {
		if err := handler(gw, session, frame); err != nil {
			fmt.Printf("Dispatcher error %v", err)
		}
	} else {
		fmt.Printf("[Dispatcher] unknown error : %v", frame.Cmd)
	}
}
