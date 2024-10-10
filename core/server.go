package core

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"kotunnel/base"
	"net"
)

type Server struct {
	openPort   int
	tunnelPort int
	secret     string
	stop       bool
}

func NewServer(openPort, tunnelPort int, secret string) *Server {
	return &Server{
		openPort:   openPort,
		tunnelPort: tunnelPort,
		secret:     secret,
		stop:       false,
	}
}

func (s *Server) Run() {
	openListener, tunnelListener, tunnelConnPool, err := s.serve()
	if err != nil {
		base.Println(31, 40, fmt.Sprintf("listener start failed: %s", err.Error()))
		return
	}

	defer func() {
		_ = openListener.Close()
		_ = tunnelListener.Close()
	}()

	// 隧道端口监听
	go func() {
		for {
			if s.stop {
				return
			}
			// 接受隧道连接
			tunnelConn, err := tunnelListener.Accept()
			if err != nil {
				_ = openListener.Close()
				base.Println(31, 40, fmt.Sprintf("port [%v] listener accept failed: %s", s.tunnelPort, err.Error()))
				return
			}

			err = s.handle(tunnelConn)
			if err != nil {
				base.Println(31, 40, fmt.Sprintf("tunnel [%v] -> [%v] create failed: %s", tunnelConn.RemoteAddr().String(), tunnelConn.LocalAddr().String(), err.Error()))
				continue
			}

			base.Println(32, 40, fmt.Sprintf("tunnel [%v] -> [%v] create success", tunnelConn.RemoteAddr().String(), tunnelConn.LocalAddr().String()))

			// 将隧道连接放入连接池
			tunnelConnPool <- tunnelConn
		}
	}()

	// 开放端口监听
	for {
		if s.stop {
			return
		}
		openConn, err := openListener.Accept()
		if err != nil {
			base.Println(31, 40, fmt.Sprintf("port [%v] listener accept failed: %s", s.openPort, err.Error()))
			return
		}
		go func() {
			err = s.copy(openConn, tunnelConnPool)
			if err != nil {
				base.Println(31, 40, fmt.Sprintf("tunnel [%v] -> [%v] connection copy fail: %s", s.tunnelPort, s.openPort, err.Error()))
				return
			}
		}()
	}
}

func (s *Server) Stop() {
	s.stop = true
}

func (s *Server) handle(conn net.Conn) (err error) {

	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()

	// 密钥验证
	var bs32 = make([]byte, 32)
	_, err = conn.Read(bs32)
	if err != nil {
		return err
	}
	// 密钥匹配
	// fmt.Println(fmt.Sprintf("%x", bs32), fmt.Sprintf("%x", sha256.Sum256([]byte(secret))))
	if fmt.Sprintf("%x", bs32) != fmt.Sprintf("%x", sha256.Sum256([]byte(s.secret))) {
		return errors.New("secret error")
	}
	// 响应验证结果
	_, err = conn.Write(base.Int64ToBytes(1, 8))
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) serve() (net.Listener, net.Listener, chan net.Conn, error) {

	var pool = make(chan net.Conn, 500)

	open, err := net.Listen("tcp", fmt.Sprintf(":%v", s.openPort))
	if err != nil {
		return nil, nil, pool, err
	}

	tunnel, err := net.Listen("tcp", fmt.Sprintf(":%v", s.tunnelPort))
	if err != nil {
		return nil, nil, pool, err
	}

	return open, tunnel, pool, nil
}

func (s *Server) copy(openConn net.Conn, tunnelConnPool chan net.Conn) error {

	defer openConn.Close()

	var tunnelConn net.Conn = nil
	for {
		tunnelConn = <-tunnelConnPool
		_, err := tunnelConn.Write(base.Int64ToBytes(1, 8))
		if err != nil {
			_ = tunnelConn.Close()
			base.Println(31, 40, fmt.Sprintf("tunnel connection write error: %s", err.Error()))
			continue
		}
		base.Println(32, 40, fmt.Sprintf("tunnel [%v] -> [%v] available", tunnelConn.RemoteAddr().String(), tunnelConn.LocalAddr().String()))
		break
	}

	defer tunnelConn.Close()

	bs8 := make([]byte, 8)
	_, err := tunnelConn.Read(bs8)
	if err != nil {
		return err
	}

	cmd := base.BytesToInt64(bs8)
	if cmd == 1 {
		base.CopyConn(tunnelConn, openConn)
	} else {
		return errors.New("bad command")
	}
	return nil
}
