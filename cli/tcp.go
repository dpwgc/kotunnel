package cli

import (
	"errors"
	"fmt"
	"kotunnel/base"
	"net"
	"time"
)

func TCP(remoteAddr string, localPort int) {
	for {
		// 连接到服务端
		remoteConn, err := net.Dial("tcp", remoteAddr)
		if err != nil {
			base.Println(31, 40, fmt.Sprintf("remote server [%v] connection failed: %s", remoteAddr, err.Error()))
			base.Logger.Error(fmt.Sprintf("error connecting to server: %v", err))
			time.Sleep(5 * time.Second)
			continue
		}

		err = tcpConn(remoteAddr, localPort, remoteConn)
		if err != nil {
			remoteConn.Close()
			base.Logger.Error(fmt.Sprintf("[%v] -> [%v] connection error: %s", localPort, remoteAddr, err.Error()))
			base.Println(33, 40, fmt.Sprintf("[%v] -> [%v] connection error: %s", localPort, remoteAddr, err.Error()))
		}

		base.Println(33, 40, fmt.Sprintf("[%v] -> [%v] connection close", localPort, remoteAddr))
	}
}

func tcpConn(remoteAddr string, localPort int, remoteConn net.Conn) error {

	var header = make([]byte, 8)
	_, err := remoteConn.Read(header)
	if err != nil {
		return err
	}

	cmd := base.BytesToInt64(header)
	if cmd == 1 {
		localConn, err := net.Dial("tcp", fmt.Sprintf(":%v", localPort))
		if err != nil {
			time.Sleep(5 * time.Second)
			return err
		}

		// 成功建立连接，return
		base.Println(32, 40, fmt.Sprintf("[%v] -> [%v] connection success", localPort, remoteAddr))
		go base.CopyConn(localConn, remoteConn)
		return nil
	}
	return errors.New("bad server command")
}
