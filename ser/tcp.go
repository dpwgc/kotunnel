package ser

import (
	"fmt"
	"kotunnel/base"
	"net"
	"sync"
	"time"
)

func TCP(openPort, clientPort int) {

	openListener, err := net.Listen("tcp", fmt.Sprintf(":%v", openPort))
	if err != nil {
		base.Logger.Error(fmt.Sprintf("error starting open listener on port %v: %v", openPort, err))
		return
	}
	defer openListener.Close()

	clientListener, err := net.Listen("tcp", fmt.Sprintf(":%v", clientPort))
	if err != nil {
		base.Logger.Error(fmt.Sprintf("error starting client listener on port %v: %v", clientPort, err))
		return
	}
	defer clientListener.Close()

	var connPool sync.Pool
	connPool.New = func() interface{} {
		return nil
	}

	go func() {
		for {
			clientConn, err := clientListener.Accept()
			if err != nil {
				base.Logger.Error(fmt.Sprintf("error accepting client connection: %v", err))
				return
			}
			connPool.Put(clientConn)
		}
	}()

	for {
		incomingConn, err := openListener.Accept()
		if err != nil {
			base.Logger.Error(fmt.Sprintf("error accepting incoming connection: %v", err))
			return
		}
		go func() {
			down := 0
			for {
				cache := connPool.Get()
				if cache != nil {
					clientConn := cache.(net.Conn)
					err = base.CopyConn(clientConn, incomingConn)
					if err != nil {
						continue
					}
				} else {
					down++
					time.Sleep(50 * time.Millisecond)
					if down > 200 {
						base.Logger.Error("failed to get connection from pool")
						incomingConn.Close()
						return
					}
				}
			}
		}()
	}
}
