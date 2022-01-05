package tcp

import (
	"fmt"
	"net"
)

type netListen struct {
	address string
	*net.TCPListener
	channelHandel func(conn *net.TCPConn)
}

func TcpCodeDe(frame []byte,network inter.INetwork) {
	data, remainingByte := utils.DecodeRound(frame)
	network.Action(func() {
		network.WriteReadChan(data)
		if len(remainingByte) > 0 {
			TcpCodeDe(remainingByte,network)
		}
	})
}
