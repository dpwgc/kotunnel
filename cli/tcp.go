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
			base.Println(31, 40, fmt.Sprintf("远程服务[%v]连接失败，5秒后重试，失败原因: %s", remoteAddr, err.Error()))
			base.Logger.Error(fmt.Sprintf("error connecting to server: %v", err))
			time.Sleep(5 * time.Second)
			continue
		}

		localConn, err := net.Dial("tcp", fmt.Sprintf(":%v", localPort))
		if err != nil {
			base.Println(31, 40, fmt.Sprintf("本地服务[%v]连接失败，5秒后重试，失败原因: %s", localPort, err.Error()))
			base.Logger.Error(fmt.Sprintf("error connecting to server: %v", err))
			time.Sleep(5 * time.Second)
			continue
		}

		go func() {
			base.CopyConn(localConn, remoteConn)
			controlChan <- true
		}()
	}
}
