package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet"
	"gnet_tool/net"
)

func main() {
	tcpServer := net.NewTcpServer("", 9000, true, true, func(frame []byte, c gnet.Conn) {
		ants.Submit(func() {
			fmt.Println(string(frame))
			if len(frame) > 0 {
				c.AsyncWrite([]byte("99998"))
				c.AsyncWrite([]byte("99997"))
				c.AsyncWrite([]byte("99999"))
			}
			return
		})
	},nil)
	tcpServer.Start()
	//time.Sleep(20 * time.Second)
}
