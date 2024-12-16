package core

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"kotunnel/base"
	"net"
)

type Client struct {
	tunnelAddr string
	localPort  int
	secret     string
	close      bool
}

func NewClient(tunnelAddr string, localPort int, secret string) *Client {
	return &Client{
		tunnelAddr: tunnelAddr,
		localPort:  localPort,
		secret:     secret,
		close:      false,
	}
}

func (c *Client) Run() {
	for {
		if c.close {
			return
		}
		// 连接到服务端
		tunnelConn, err := net.Dial("tcp", c.tunnelAddr)
		if err != nil {
			base.Tips(base.Red, fmt.Sprintf("tunnel server [%v] connection failed: %s", c.tunnelAddr, err.Error()), 5)
			continue
		}
		// 建立隧道（与服务端建立长连接）
		err = c.hangOn(tunnelConn)
		if err != nil {
			base.Tips(base.Red, fmt.Sprintf("tunnel [%v] -> [%v] create failed: %s", tunnelConn.LocalAddr().String(), c.tunnelAddr, err.Error()), 5)
		} else {
			base.Tips(base.Green, fmt.Sprintf("tunnel [%v] -> [%v] create success", tunnelConn.LocalAddr().String(), c.tunnelAddr))
		}
	}
}

func (c *Client) Close() {
	c.close = true
}

func (c *Client) hangOn(tunnelConn net.Conn) (err error) {

	defer func() {
		if err != nil {
			_ = tunnelConn.Close()
		}
	}()

	// 密钥验证
	bs32 := sha256.Sum256([]byte(c.secret))
	_, err = tunnelConn.Write(bs32[:32])
	if err != nil {
		return err
	}

	var bs8 = make([]byte, 8)

	// 服务端检查客户端是否通过验证（不断开就算通过）
	_, err = tunnelConn.Read(bs8)
	if err != nil {
		return err
	}

	// 服务端是否要使用这个连接
	_, err = tunnelConn.Read(bs8)
	if err != nil {
		return err
	}

	if base.Valid(bs8) {
		// 连接服务端
		localConn, err := net.Dial("tcp", fmt.Sprintf(":%v", c.localPort))
		if err != nil {
			return err
		}
		// 回应服务端
		_, err = tunnelConn.Write(base.Ping())
		if err != nil {
			return err
		}
		// 成功建立连接，return
		go base.Copy(localConn, tunnelConn)
		return nil
	}
	return errors.New("bad command")
}
