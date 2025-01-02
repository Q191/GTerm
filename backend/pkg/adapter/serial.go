package adapter

import (
	"encoding/json"
	"github.com/MisakaTAT/GTerm/backend/pkg/terminal"
	"github.com/MisakaTAT/GTerm/backend/pkg/types"
	"github.com/gorilla/websocket"
	"go.bug.st/serial"
	"go.uber.org/zap"
)

type Serial struct {
	port   serial.Port
	ws     *websocket.Conn
	logger *zap.Logger
}

func NewSerial(ws *websocket.Conn, logger *zap.Logger) *Serial {
	return &Serial{
		ws:     ws,
		logger: logger,
	}
}

func (s *Serial) SerialPorts() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		panic(err)
	}
	return ports
}

func (s *Serial) Open(portName string) error {
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		return err
	}
	s.port = port
	return nil
}

func (s *Serial) close() {
	if s.port != nil {
		_ = s.port.Close()
	}
}

func (s *Serial) Input(quitSignal chan bool) {
	defer s.setQuit(quitSignal)
	for {
		select {
		case <-quitSignal:
			return
		default:
			_, data, err := s.ws.ReadMessage()
			if err != nil {
				s.logger.Error("read message failed", zap.Error(err))
				return
			}
			msg := &terminal.Payload{}
			_ = json.Unmarshal(data, &msg)
			if msg.Type == terminal.Command {
				if _, err = s.port.Write([]byte(msg.Cmd)); err != nil {
					s.logger.Error("write to serial port failed", zap.Error(err))
				}
			}
		}
	}
}

func (s *Serial) Output(quitSignal chan bool) {
	defer s.setQuit(quitSignal)
	if s.port == nil {
		s.logger.Error("serial port is closed")
		s.setQuit(quitSignal)
	}
	buff := make([]byte, 100)
	for {
		select {
		case <-quitSignal:
			return
		default:
			for {
				n, err := s.port.Read(buff)
				if err != nil {
					s.logger.Error("error reading from serial port", zap.Error(err))
					return
				}
				if n == 0 || len(buff[:n]) == 0 {
					continue
				}
				if err = s.ws.WriteJSON(&types.Message{
					Type:    types.MessageTypeData,
					Content: string(buff[:n]),
				}); err != nil {
					s.logger.Error("write message failed", zap.Error(err))
				}
			}
		}
	}
}

func (s *Serial) Wait(quitSignal chan bool) {
	defer s.setQuit(quitSignal)
	<-quitSignal
	s.close()
}

func (s *Serial) setQuit(ch chan bool) {
	ch <- true
}
