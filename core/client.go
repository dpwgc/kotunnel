package core

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"kotunnel/base"
	"net"
	"time"
)

type Client struct {
	tunnelAddr string
	localPort  int
	secret     string
	stop       bool
}

func NewClient(tunnelAddr string, localPort int, secret string) *Client {
	return &Client{
		tunnelAddr: tunnelAddr,
		localPort:  localPort,
		secret:     secret,
		stop:       false,
	}
}

func (c *Client) Run() {
	for {
		if c.stop {
			return
		}
		// 连接到服务端
		tunnelConn, err := net.Dial("tcp", c.tunnelAddr)
		if err != nil {
			base.Println(31, 40, fmt.Sprintf("tunnel server [%v] connection failed: %s", c.tunnelAddr, err.Error()))
			time.Sleep(5 * time.Second)
			continue
		}

		// 建立隧道
		err = c.handle(tunnelConn)
		if err != nil {
			base.Println(31, 40, fmt.Sprintf("tunnel [%v] -> [%v] create failed: %s", tunnelConn.LocalAddr().String(), c.tunnelAddr, err.Error()))
		} else {
			base.Println(32, 40, fmt.Sprintf("tunnel [%v] -> [%v] create success", tunnelConn.LocalAddr().String(), c.tunnelAddr))
		}
	}
}

func (c *Client) Stop() {
	c.stop = true
}

func (c *Client) handle(tunnelConn net.Conn) (err error) {

	useSleep := false

	defer func() {
		if err != nil {
			_ = tunnelConn.Close()
			if useSleep {
				time.Sleep(5 * time.Second)
			}
		}
	}()

	// 密钥验证
	bs32 := sha256.Sum256([]byte(c.secret))
	_, err = tunnelConn.Write(bs32[:32])
	if err != nil {
		useSleep = true
		return err
	}

	var bs8 = make([]byte, 8)

	// 是否通过验证（不断开就算通过）
	_, err = tunnelConn.Read(bs8)
	if err != nil {
		useSleep = true
		return err
	}

	_, err = tunnelConn.Read(bs8)
	if err != nil {
		return err
	}

	cmd := base.BytesToInt64(bs8)
	if cmd == 1 {

		localConn, err := net.Dial("tcp", fmt.Sprintf(":%v", c.localPort))
		if err != nil {
			useSleep = true
			return err
		}

		_, err = tunnelConn.Write(base.Int64ToBytes(1, 8))
		if err != nil {
			return err
		}

		// 成功建立连接，return
		go base.CopyConn(localConn, tunnelConn)
		return nil
	}
	return errors.New("bad command")
}
