package ser

import (
	"fmt"
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

	clientConnChan := make(chan net.Conn, 100)

	go func() {
		for {
			clientConn, err := clientListener.Accept()
			if err != nil {
				base.Logger.Error(fmt.Sprintf("error accepting client connection: %v", err))
				return
			}
			clientConnChan <- clientConn
			fmt.Println("push clientConn")
		}
	}()

	for {
		incomingConn, err := listener.Accept()
		if err != nil {
			base.Logger.Error(fmt.Sprintf("error accepting incoming connection: %v", err))
			return
		}
		clientConn := <-clientConnChan
		fmt.Println("pop clientConn")
		go base.CopyConn(clientConn, incomingConn)
	}
}
