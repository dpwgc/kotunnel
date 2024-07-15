package ser

import (
	"fmt"
	"io"
	"kotunnel/base"
	"net"
)

func TCP(listenPort, clientPort int) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", listenPort))
	if err != nil {
		base.Logger.Error(fmt.Sprintf("error starting listener on port %v: %v", listenPort, err))
		return
	}
	defer listener.Close()

	clientListener, err := net.Listen("tcp", fmt.Sprintf(":%v", clientPort))
	if err != nil {
		base.Logger.Error(fmt.Sprintf("error starting listener on port %v: %v", clientPort, err))
		return
	}
	defer clientListener.Close()

	for {
		clientConn, err := clientListener.Accept()
		if err != nil {
			base.Logger.Error(fmt.Sprintf("error accepting client connection: %v", err))
			continue
		}

		go handleClient(clientConn, listener)
	}
}

func handleClient(clientConn net.Conn, listener net.Listener) {
	defer clientConn.Close()

	for {
		incomingConn, err := listener.Accept()
		if err != nil {
			base.Logger.Error(fmt.Sprintf("error accepting incoming connection: %v", err))
			continue
		}

		go handleConnection(clientConn, incomingConn)
	}
}

func handleConnection(clientConn, incomingConn net.Conn) {
	defer clientConn.Close()
	defer incomingConn.Close()
	go io.Copy(clientConn, incomingConn)
	io.Copy(incomingConn, clientConn)
}
