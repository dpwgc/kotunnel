package ser

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"kotunnel/base"
	"net"
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
				base.Logger.Error(fmt.Sprintf("tunnel listener error: %s", err.Error()))
				return
			}

			err = tcpHandle(tunnelConn, secret)
			if err != nil {
				base.Logger.Error(fmt.Sprintf("tunnel [%s] -> [%v] connection fail: %s", tunnelConn.RemoteAddr().String(), tunnelPort, err.Error()))
				continue
			}

			// 将隧道连接放入连接池
			tunnelConnPool <- tunnelConn
		}
	}()

	// 开放端口监听
	for {
		openConn, err := openListener.Accept()
		if err != nil {
			base.Logger.Error(fmt.Sprintf("open listener error: %s", err.Error()))
			return
		}
		go func() {
			err = tcpCopy(openConn, tunnelConnPool)
			if err != nil {
				base.Logger.Error(fmt.Sprintf("tcp connection copy fail: %s", err.Error()))
				return
			}
		}()
	}
}

func tcpHandle(conn net.Conn, secret string) (err error) {

	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	// 密钥验证
	var bs = make([]byte, 32)
	_, err = conn.Read(bs)
	if err != nil {
		return err
	}
	// 密钥匹配
	if fmt.Sprintf("%x", bs) == fmt.Sprintf("%x", sha256.Sum256([]byte(secret))) {
		return errors.New("secret error")
	}
	// 响应验证结果
	_, err = conn.Write(base.Int64ToBytes(1, 8))
	if err != nil {
		return err
	}
	return nil
}

func tcpServe(openPort, tunnelPort int) (net.Listener, net.Listener, chan net.Conn, error) {

	var pool = make(chan net.Conn, 500)

	open, err := net.Listen("tcp", fmt.Sprintf(":%v", openPort))
	if err != nil {
		return nil, nil, pool, err
	}

	tunnel, err := net.Listen("tcp", fmt.Sprintf(":%v", tunnelPort))
	if err != nil {
		return nil, nil, pool, err
	}

	return open, tunnel, pool, nil
}

func tcpCopy(openConn net.Conn, tunnelConnPool chan net.Conn) error {

	defer openConn.Close()

	var tunnelConn net.Conn = nil
	for {
		tunnelConn = <-tunnelConnPool
		_, err := tunnelConn.Write(base.Int64ToBytes(1, 8))
		if err != nil {
			tunnelConn.Close()
			base.Logger.Error(fmt.Sprintf("tunnel connection write error: %s", err.Error()))
			continue
		}
		break
	}

	defer tunnelConn.Close()

	bs8 := make([]byte, 8)
	_, err := tunnelConn.Read(bs8)
	if err != nil {
		return err
	}

	cmd := base.BytesToInt64(bs8)
	if cmd == 1 {
		base.CopyConn(tunnelConn, openConn)
	} else {
		return errors.New("bad command")
	}
	return nil
}
