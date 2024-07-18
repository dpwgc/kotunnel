package ser

import (
	"crypto/sha256"
	"fmt"
	"kotunnel/base"
	"net"
	"sync"
	"time"
)

func TCP(openPort, tunnelPort int, secret string) {

	openListener, tunnelListener, tunnelConnPool, err := tcpServe(openPort, tunnelPort)
	if err != nil {
		base.Logger.Error(fmt.Sprintf("start listener error: %s", err.Error()))
		return
	}

	defer func() {
		openListener.Close()
		tunnelListener.Close()
	}()

	// 隧道端口监听
	go func() {
		for {
			// 接受隧道连接
			tunnelConn, err := tunnelListener.Accept()
			if err != nil {
				openListener.Close()
				base.Logger.Error(fmt.Sprintf("tunnel connection accept error: %s", err.Error()))
				return
			}

			// 密钥验证
			var bs = make([]byte, 32)
			_, err = tunnelConn.Read(bs)
			if err != nil {
				tunnelConn.Close()
				base.Logger.Error(fmt.Sprintf("tunnel connection read error: %s", err.Error()))
				continue
			}
			// 密钥匹配
			if fmt.Sprintf("%x", bs) == fmt.Sprintf("%x", sha256.Sum256([]byte(secret))) {
				tunnelConn.Close()
				base.Logger.Error("tunnel connection secret error")
				continue
			}
			// 响应验证结果
			_, err = tunnelConn.Write(base.Int64ToBytes(1, 8))
			if err != nil {
				tunnelConn.Close()
				base.Logger.Error(fmt.Sprintf("tunnel connection write error: %s", err.Error()))
				continue
			}

			// 将隧道连接放入连接池
			tunnelConnPool.Put(tunnelConn)
		}
	}()

	// 开放端口监听
	for {
		openConn, err := openListener.Accept()
		if err != nil {
			base.Logger.Error(fmt.Sprintf("open connection accept error: %s", err.Error()))
			return
		}
		go tcpHandle(openConn, tunnelConnPool)
	}
}

func tcpServe(openPort, tunnelPort int) (net.Listener, net.Listener, *sync.Pool, error) {

	var pool sync.Pool
	pool.New = func() interface{} {
		return nil
	}

	open, err := net.Listen("tcp", fmt.Sprintf(":%v", openPort))
	if err != nil {
		return nil, nil, &pool, err
	}

	tunnel, err := net.Listen("tcp", fmt.Sprintf(":%v", tunnelPort))
	if err != nil {
		return nil, nil, &pool, err
	}

	return open, tunnel, &pool, nil
}

func tcpHandle(openConn net.Conn, tunnelConnPool *sync.Pool) {
	retry := 0
	for {
		retry++
		// 超过最大重试次数
		if retry > 200 {
			openConn.Close()
			base.Logger.Error("reached maximum retry limit")
			return
		}

		cache := tunnelConnPool.Get()
		if cache == nil {
			base.Logger.Error("connection pool is empty")
			time.Sleep(50 * time.Millisecond)
			continue
		}

		tunnelConn := cache.(net.Conn)
		_, err := tunnelConn.Write(base.Int64ToBytes(1, 8))
		if err != nil {
			tunnelConn.Close()
			base.Logger.Error(fmt.Sprintf("tunnel connection accept error: %s", err.Error()))
			time.Sleep(50 * time.Millisecond)
			continue
		}

		base.CopyConn(tunnelConn, openConn)
		return
	}
}
