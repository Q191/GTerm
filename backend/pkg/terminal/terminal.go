package terminal

import (
	"github.com/gorilla/websocket"
	"sync"
)

const (
	Resize  = "resize"
	Command = "cmd"
)

type Handler interface {
	Input(quitSignal chan bool)
	Output(quitSignal chan bool)
	Wait(quitSignal chan bool)
}

type Config struct {
	WebSocket *websocket.Conn
}

type Payload struct {
	Type string `json:"type"`
	Cmd  string `json:"cmd"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}

type Terminal struct {
	handler  Handler
	conf     *Config
	stopFunc func(ws *websocket.Conn)
}

func NewTerminal(conf *Config, handler Handler, stopFunc func(ws *websocket.Conn)) *Terminal {
	return &Terminal{
		conf:     conf,
		handler:  handler,
		stopFunc: stopFunc,
	}
}

func (t *Terminal) Start() {
	var wg sync.WaitGroup
	wg.Add(3)

	quitSignal := make(chan bool, 3)
	outputClosed := make(chan struct{})

	go func() {
		defer wg.Done()
		t.handler.Input(quitSignal)
	}()

	go func() {
		defer wg.Done()
		t.handler.Output(quitSignal)
		close(outputClosed)
	}()

	go func() {
		defer wg.Done()
		t.handler.Wait(quitSignal)
	}()

	go func() {
		// ensure flush write buffer before close
		// must output handler close before with websocket close
		<-outputClosed
		t.stopFunc(t.conf.WebSocket)
	}()

	wg.Wait()
}
