package main

import "github.com/Byfengfeng/gnet_tool/net"

func main() {
	tcpServer := net.NewTcpServer("tcp6","", 9000, true)
	err := tcpServer.Start()
	if err != nil {
		panic(err)
	}
	//strs := make(chan int)
	//go func() {
	//	for  {
	//		a := <-strs
	//		if a == 0 {
	//			break
	//		}
	//		fmt.Println(a)
	//	}
	//
	//	return
	//}()
	//
	//for i:= 1; i < 5; i++ {
	//	strs <- i
	//	time.Sleep(time.Second)
	//}
	////strs <- 0
	//fmt.Println("close chan")
	//close(strs)

	<-make(chan struct{})
}