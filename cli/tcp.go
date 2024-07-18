package cli

import (
	"errors"
	"fmt"
	"kotunnel/base"
	"net"
	"time"
)

func TCP(tunnelAddr string, localPort int) {
	for {
		// 连接到服务端
		tunnelConn, err := net.Dial("tcp", tunnelAddr)
		if err != nil {
			base.Println(31, 40, fmt.Sprintf("tunnel server [%v] connection failed: %s", tunnelAddr, err.Error()))
			base.Logger.Error(fmt.Sprintf("error connecting to tunnel server: %v", err))
			time.Sleep(5 * time.Second)
			continue
		}

		err = tcpHandle(localPort, tunnelConn)
		if err != nil {
			base.Logger.Error(fmt.Sprintf("tunnel [%v] -> [%v] create failed: %s", localPort, tunnelAddr, err.Error()))
			base.Println(33, 40, fmt.Sprintf("tunnel [%v] -> [%v] create failed: %s", localPort, tunnelAddr, err.Error()))
		} else {
			base.Println(32, 40, fmt.Sprintf("tunnel [%v] -> [%v] create success", localPort, tunnelAddr))
		}
	}
}

func tcpHandle(localPort int, tunnelConn net.Conn) (err error) {

	defer func() {
		if err != nil {
			tunnelConn.Close()
		}
	}()

	var header = make([]byte, 8)
	_, err = tunnelConn.Read(header)
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
		go base.CopyConn(localConn, tunnelConn)
		return nil
	}
	return errors.New("bad server command")
}
