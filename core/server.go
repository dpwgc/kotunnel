package core

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"kotunnel/base"
	"net"
)

type Server struct {
	name           string
	openPort       int
	tunnelPort     int
	tunnelPool     chan net.Conn
	openListener   net.Listener
	tunnelListener net.Listener
	secret         string
}

func NewServer(name string, openPort, tunnelPort int, maxConn int, secret string) *Server {
	if maxConn <= 0 {
		maxConn = 1000
	}
	return &Server{
		name:       name,
		openPort:   openPort,
		tunnelPort: tunnelPort,
		tunnelPool: make(chan net.Conn, maxConn),
		secret:     secret,
	}
}

func (c *Server) Name() string {
	return fmt.Sprintf("{server:%s}", c.name)
}

func (s *Server) Run() (oErr, tErr error) {
	base.Info(fmt.Sprintf("%s server [%v] -> [%v] launch", s.Name(), s.tunnelPort, s.openPort))
	defer s.Close()
	go func() {
		defer s.Close()
		oErr = s.listenOpenPort()
	}()
	tErr = s.listenTunnelPort()
	return oErr, tErr
}

func (s *Server) listenOpenPort() (err error) {

	s.openListener, err = net.Listen("tcp", fmt.Sprintf(":%v", s.openPort))
	if err != nil {
		return err
	}
	defer s.Close()

	// 开放端口监听
	for {
		conn, err := s.openListener.Accept()
		if err != nil {
			return err
		}
		go func() {
			err = s.copy(conn)
			if err != nil {
				base.Error(fmt.Sprintf("%s tunnel [%v] -> [%v] connection copy fail: %s", s.Name(), s.tunnelPort, s.openPort, err.Error()))
			} else {
				base.Success(fmt.Sprintf("%s tunnel [%v] -> [%v] connection copy success", s.Name(), s.tunnelPort, s.openPort))
			}
		}()
	}
}

func (s *Server) listenTunnelPort() (err error) {

	s.tunnelListener, err = net.Listen("tcp", fmt.Sprintf(":%v", s.tunnelPort))
	if err != nil {
		return err
	}
	defer s.Close()

	for {
		// 接受隧道连接
		conn, err := s.tunnelListener.Accept()
		if err != nil {
			_ = s.openListener.Close()
			return err
		}
		// 验证密钥，验证通过后，将连接挂在服务端上，保持长连接
		err = s.hangOn(conn)
		if err != nil {
			base.Error(fmt.Sprintf("%s tunnel [%v] -> [%v] create failed: %s", s.Name(), conn.RemoteAddr().String(), conn.LocalAddr().String(), err.Error()))
			continue
		}
		base.Success(fmt.Sprintf("%s tunnel [%v] -> [%v] create success", s.Name(), conn.RemoteAddr().String(), conn.LocalAddr().String()))
		// 将隧道连接放入连接池
		s.tunnelPool <- conn
	}
}

func (s *Server) Close() {
	defer func() {
		_ = recover()
	}()
	_ = s.openListener.Close()
	_ = s.tunnelListener.Close()
	close(s.tunnelPool)
}

func (s *Server) hangOn(conn net.Conn) (err error) {

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
	if len(s.secret) > 0 && fmt.Sprintf("%x", bs32) != fmt.Sprintf("%x", sha256.Sum256([]byte(s.secret))) {
		return errors.New("secret error")
	}
	// 响应验证结果
	_, err = conn.Write(base.Ping())
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) copy(openConn net.Conn) error {

	defer openConn.Close()

	var tunnelConn net.Conn = nil
	for {
		tunnelConn = <-s.tunnelPool
		// 向客户端发送ping命令，验证这个连接是否存活
		_, err := tunnelConn.Write(base.Ping())
		if err != nil {
			_ = tunnelConn.Close()
			base.Error(fmt.Sprintf("%s tunnel connection write error: %s", s.Name(), err.Error()))
			continue
		}
		base.Success(fmt.Sprintf("%s tunnel [%v] -> [%v] available", s.Name(), tunnelConn.RemoteAddr().String(), tunnelConn.LocalAddr().String()))
		break
	}

	defer tunnelConn.Close()

	// 获取客户端回应
	bs8 := make([]byte, 8)
	_, err := tunnelConn.Read(bs8)
	if err != nil {
		return err
	}

	// 如果客户端成功回应，说明这个连接是可用的，开始交换复制连接
	if base.Valid(bs8) {
		base.Copy(tunnelConn, openConn)
		return nil
	}
	return errors.New("bad command (check to see if the port you set is occupied by another app)")
}
