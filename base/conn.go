package base

import (
	"io"
	"net"
	"sync"
)

func CopyConn(local, remote net.Conn) error {

	var err error = nil

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		_, er := io.Copy(remote, local)
		if er != nil && er != io.EOF {
			err = er
		}
		wg.Done()
	}()
	go func() {
		_, er := io.Copy(local, remote)
		if er != nil && er != io.EOF {
			err = er
		}
		wg.Done()
	}()
	wg.Wait()
	return err
}
