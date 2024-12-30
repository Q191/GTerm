package initialize

import (
	"net"
	"net/http"
)

type HTTPListenerPort int

func InitHTTPServer() *HTTPListenerPort {
	// if address is ":0" then auto assign a port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	go func() {
		if err = http.Serve(listener, nil); err != nil {
			panic(err)
		}
	}()
	port := HTTPListenerPort(listener.Addr().(*net.TCPAddr).Port)
	return &port
}
