package cli

import (
	"fmt"
	"kotunnel/base"
	"net"
	"time"
)

func TCP(remoteAddr string, localPort int) {
	controlChan := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		controlChan <- true
	}
	for {
		select {
		case <-controlChan:
			fmt.Printf("\033[1;36;40m%s\033[0m\n", fmt.Sprintf("尝试连接远程服务[%v]......", remoteAddr))
			// 连接到服务端
			remoteConn, err := net.Dial("tcp", remoteAddr)
			if err != nil {
				fmt.Printf("\033[1;31;40m%s\033[0m\n", fmt.Sprintf("远程服务[%v]连接失败，5秒后重试，失败原因: %s", remoteAddr, err.Error()))
				base.Logger.Error(fmt.Sprintf("error connecting to server: %v", err))
				time.Sleep(5 * time.Second)
				continue
			}

			fmt.Printf("\033[1;32;40m%s\033[0m\n", fmt.Sprintf("远程服务[%v]连接成功", remoteAddr))

			fmt.Printf("\033[1;36;40m%s\033[0m\n", fmt.Sprintf("尝试连接本地服务[%v]......", localPort))
			localConn, err := net.Dial("tcp", fmt.Sprintf(":%v", localPort))
			if err != nil {
				fmt.Printf("\033[1;31;40m%s\033[0m\n", fmt.Sprintf("本地服务[%v]连接失败，5秒后重试，失败原因: %s", localPort, err.Error()))
				base.Logger.Error(fmt.Sprintf("error connecting to server: %v", err))
				time.Sleep(5 * time.Second)
				continue
			}

			fmt.Printf("\033[1;32;40m%s\033[0m\n", fmt.Sprintf("本地服务[%v]连接成功", localPort))

			go func() {
				base.CopyConn(localConn, remoteConn)
				controlChan <- true
			}()
		}
	}
}
