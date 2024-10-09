package base

import (
	"encoding/binary"
	"io"
	"net"
)

func CopyConn(local, remote net.Conn) {
	go func() {
		_, _ = io.Copy(local, remote)
		_ = remote.Close()
		_ = local.Close()
	}()
	_, _ = io.Copy(remote, local)
	_ = remote.Close()
	_ = local.Close()
}

func Int64ToBytes(num int64, len int) []byte {
	byteArray := make([]byte, len)
	binary.LittleEndian.PutUint64(byteArray, uint64(num))
	return byteArray
}

func BytesToInt64(bytes []byte) int64 {
	return int64(binary.LittleEndian.Uint64(bytes[:]))
}
