package cli

import (
	"fmt"
	"kotunnel/base"
	"net"
	"time"
)

func TCP(remoteAddr string, localPort int) {
	for {
		time.Sleep(1 * time.Second)
		fmt.Printf("\033[1;36;40m%s\033[0m\n", "尝试连接远程服务......")
		// 连接到服务端
		remoteConn, err := net.Dial("tcp", remoteAddr)
		if err != nil {
			fmt.Printf("\033[1;31;40m%s\033[0m\n", fmt.Sprintf("远程服务连接失败，5秒后重试，失败原因: %s", err.Error()))
			base.Logger.Error(fmt.Sprintf("error connecting to server: %v", err))
			time.Sleep(4 * time.Second)
			continue
		}

		fmt.Printf("\033[1;32;40m%s\033[0m\n", "远程服务连接成功")

		localConn, err := net.Dial("tcp", fmt.Sprintf(":%v", localPort))
		if err != nil {
			fmt.Printf("\033[1;31;40m%s\033[0m\n", fmt.Sprintf("本地服务连接失败，5秒后重试，失败原因: %s", err.Error()))
			base.Logger.Error(fmt.Sprintf("error connecting to server: %v", err))
			time.Sleep(4 * time.Second)
			continue
		}

		fmt.Printf("\033[1;32;40m%s\033[0m\n", "本地服务连接成功")

		base.CopyConn(localConn, remoteConn)
	}
}
