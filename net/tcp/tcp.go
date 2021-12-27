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

func NewNetListen(addr string,channelHandel func(con *net.TCPConn)) *netListen {
	return &netListen{address: addr,channelHandel: channelHandel}
}

func (n * netListen) Start() error {
	addr, err := net.ResolveTCPAddr("tcp", n.address)
	if err != nil {
		return err
	}
	n.TCPListener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("listener", fmt.Errorf("%s", err))
			}
		}()
		for {
			tcpConn, err := n.TCPListener.AcceptTCP()
			if err != nil {
				fmt.Println("client channel exit", fmt.Errorf("%s", err))
			}
			n.channelHandel(tcpConn)
		}
	}()
	fmt.Println("tcp listen start success")
	return nil
}
