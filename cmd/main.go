package main

import (
	"fmt"
	"github.com/Byfengfeng/gnet_tool/net"
	"runtime/debug"
	"time"
)

func main() {
	//go freeOs()
	tcpServer := net.NewTcpServer("tcp6","", 9000, true)
	err := tcpServer.Start()
	if err != nil {
		panic(err)
	}

	//strs := make(chan []byte)
	//go func() {
	//	for  {
	//		a := <-strs
	//		if len(a) == 0 {
	//			break
	//		}
	//		fmt.Println(a)
	//	}
	//
	//	return
	//}()
	//
	//for i:= 1; i < 50000; i++ {
	//	strs <- []byte("abcdefg123")
	//}
	////strs <- 0
	//time.Sleep(5 * time.Second)
	//fmt.Println("close chan")
	//close(strs)
	//debug.FreeOSMemory()
	//time.Sleep(10 * time.Second)

	<-make(chan struct{})
}

func freeOs()  {
	timer := time.NewTimer(3 * time.Second)
		select {
		case <- timer.C:
			fmt.Println("FreeOSMemory")
			debug.FreeOSMemory()
			freeOs()
		}
}