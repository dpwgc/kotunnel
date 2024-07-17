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

	var clientConnPool sync.Pool
	clientConnPool.New = func() interface{} {
		return nil
	}

	go func() {
		for {
			clientConn, err := clientListener.Accept()
			if err != nil {
				base.Logger.Error(fmt.Sprintf("error accepting client connection: %v", err))
				return
			}
			clientConnPool.Put(clientConn)
		}
	}()

	for {
		openConn, err := openListener.Accept()
		if err != nil {
			base.Logger.Error(fmt.Sprintf("error accepting open connection: %v", err))
			return
		}
		go func() {
			down := 0
			for {
				cache := clientConnPool.Get()
				if cache != nil {
					clientConn := cache.(net.Conn)
					_, err = clientConn.Write(base.Int64ToBytes(1, 8))
					if err != nil {
						base.Logger.Error(fmt.Sprintf("error writing client connection: %v", err))
						clientConn.Close()
						time.Sleep(50 * time.Millisecond)
						continue
					}
					base.CopyConn(clientConn, openConn)
					return
				} else {
					down++
					if down > 200 {
						base.Logger.Error("failed to get connection from pool")
						openConn.Close()
						return
					}
					time.Sleep(50 * time.Millisecond)
				}
			}
		}()
	}
}
