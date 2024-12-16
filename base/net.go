package base

import (
	"encoding/binary"
	"io"
	"net"
)

func Copy(local, remote net.Conn) {
	go func() {
		_, _ = io.Copy(local, remote)
		_ = remote.Close()
		_ = local.Close()
	}()
	_, _ = io.Copy(remote, local)
	_ = remote.Close()
	_ = local.Close()
}

func Ping() []byte {
	byteArray := make([]byte, 8)
	binary.LittleEndian.PutUint64(byteArray, uint64(1))
	return byteArray
}

func Valid(bytes []byte) bool {
	cmd := int64(binary.LittleEndian.Uint64(bytes[:]))
	return cmd == 1
}
