package cli

import (
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
			time.Sleep(1 * time.Second)
			continue
		}

		base.Println(32, 40, fmt.Sprintf("remote server [%v] connection success", remoteAddr))

		func() {

			var header = make([]byte, 8)
			_, err = remoteConn.Read(header)
			if err != nil {
				remoteConn.Close()
				base.Println(33, 40, fmt.Sprintf("remote server [%v] connection close: %s", remoteAddr, err.Error()))
				base.Logger.Error(fmt.Sprintf("error reading remote connection: %v", err))
				return
			}

			cmd := base.BytesToInt64(header)
			if cmd == 1 {
				localConn, err := net.Dial("tcp", fmt.Sprintf(":%v", localPort))
				if err != nil {
					remoteConn.Close()
					base.Println(31, 40, fmt.Sprintf("local server [%v] connection failed: %s", localPort, err.Error()))
					base.Logger.Error(fmt.Sprintf("error connecting to server: %v", err))
					time.Sleep(1 * time.Second)
					return
				}
				base.Println(32, 40, fmt.Sprintf("local server [%v] connection success", localPort))

				go base.CopyConn(localConn, remoteConn)

			} else {
				remoteConn.Close()
				base.Println(31, 40, fmt.Sprintf("remote server [%v] command failed: %s", localPort, err.Error()))
				base.Logger.Error(fmt.Sprintf("error command from server: %v", err))
				time.Sleep(1 * time.Second)
			}
		}()
	}
}
