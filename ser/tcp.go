package ser

import (
	"fmt"
	"kotunnel/base"
	"net"
	"sync"
	"time"
)

func TCP(openPort, clientPort int) {

	openListener, clientListener, clientConnPool, err := tcpServe(openPort, clientPort)
	if err != nil {
		base.Logger.Error(fmt.Sprintf("error starting listener on port %v: %v", openPort, err))
		return
	}

	defer func() {
		openListener.Close()
		clientListener.Close()
	}()

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
		go tcpHandle(openConn, clientConnPool)
	}
}

func tcpServe(openPort, clientPort int) (net.Listener, net.Listener, *sync.Pool, error) {

	var pool sync.Pool
	pool.New = func() interface{} {
		return nil
	}

	openListener, err := net.Listen("tcp", fmt.Sprintf(":%v", openPort))
	if err != nil {
		return nil, nil, &pool, err
	}

	clientListener, err := net.Listen("tcp", fmt.Sprintf(":%v", clientPort))
	if err != nil {
		return nil, nil, &pool, err
	}

	return openListener, clientListener, &pool, nil
}

func tcpHandle(openConn net.Conn, clientConnPool *sync.Pool) {
	retry := 0
	for {
		retry++
		// 超过最大重试次数
		if retry > 200 {
			openConn.Close()
			base.Logger.Error("retry > 200")
			return
		}

		cache := clientConnPool.Get()
		if cache == nil {
			base.Logger.Error("failed to get connection from pool")
			time.Sleep(50 * time.Millisecond)
			continue
		}

		clientConn := cache.(net.Conn)
		_, err := clientConn.Write(base.Int64ToBytes(1, 8))
		if err != nil {
			clientConn.Close()
			base.Logger.Error(fmt.Sprintf("error writing client connection: %v", err))
			time.Sleep(50 * time.Millisecond)
			continue
		}

		base.CopyConn(clientConn, openConn)
		return
	}
}
