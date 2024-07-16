package base

import (
	"io"
	"net"
	"sync"
)

func CopyConn(a, b net.Conn) error {

	var err error = nil

	defer a.Close()
	defer b.Close()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		_, err = io.Copy(b, a)
		a.Close()
		b.Close()
		wg.Done()
	}()
	go func() {
		_, err = io.Copy(a, b)
		a.Close()
		b.Close()
		wg.Done()
	}()
	wg.Wait()
	return err
}
