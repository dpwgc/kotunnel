package cli

import (
	"fmt"
	"kotunnel/base"
	"net"
	"time"
)

func TCP(remoteAddr string, localPort int, tunnelNum int) {
	controlChan := make(chan bool, tunnelNum)
	for i := 0; i < 100; i++ {
		controlChan <- true
	}
	for {
		_ = <-controlChan
		// 连接到服务端
		remoteConn, err := net.Dial("tcp", remoteAddr)
		if err != nil {
			controlChan <- true
			base.Println(31, 40, fmt.Sprintf("remote server [%v] connection failed: %s", remoteAddr, err.Error()))
			base.Logger.Error(fmt.Sprintf("error connecting to server: %v", err))
			time.Sleep(1 * time.Second)
			continue
		}

		localConn, err := net.Dial("tcp", fmt.Sprintf(":%v", localPort))
		if err != nil {
			remoteConn.Close()
			controlChan <- true
			base.Println(31, 40, fmt.Sprintf("local server [%v] connection failed: %s", localPort, err.Error()))
			base.Logger.Error(fmt.Sprintf("error connecting to server: %v", err))
			time.Sleep(1 * time.Second)
			continue
		}

		go func() {
			base.CopyConn(localConn, remoteConn)
			localConn.Close()
			remoteConn.Close()
			controlChan <- true
		}()
	}
}
