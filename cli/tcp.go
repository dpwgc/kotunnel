package cli

import (
	"fmt"
	"io"
	"kotunnel/base"
	"net"
	"time"
)

func TCP(remoteAddr string, localPort int) {
	for {
		time.Sleep(1 * time.Second)
		fmt.Printf("\033[1;36;40m%s\033[0m\n", "尝试连接服务端......")
		// 连接到服务端
		serverConn, err := net.Dial("tcp", remoteAddr)
		if err != nil {
			fmt.Printf("\033[1;31;40m%s\033[0m\n", fmt.Sprintf("连接失败，5秒后重试，失败原因: %s", err.Error()))
			base.Logger.Error(fmt.Sprintf("error connecting to server: %v", err))
			time.Sleep(4 * time.Second)
			continue
		}

		localListener, err := net.Listen("tcp", fmt.Sprintf(":%v", localPort))
		if err != nil {
			base.Logger.Error(fmt.Sprintf("error starting local listener on port %v: %v", localPort, err))
			return
		}

		for {
			localConn, err := localListener.Accept()
			if err != nil {
				base.Logger.Error(fmt.Sprintf("error accepting local connection: %v", err))
				continue
			}

			go handleConnection(localConn, serverConn)
		}
	}
}

func handleConnection(localConn, serverConn net.Conn) {
	defer localConn.Close()
	go io.Copy(serverConn, localConn)
	io.Copy(localConn, serverConn)
}
