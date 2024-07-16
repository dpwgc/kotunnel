package base

import (
	"io"
	"net"
	"sync"
)

func CopyConn(a, b net.Conn) {

	defer a.Close()
	defer b.Close()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		_, _ = io.Copy(b, a)
		wg.Done()
	}()
	go func() {
		_, _ = io.Copy(a, b)
		wg.Done()
	}()
	wg.Wait()
}
