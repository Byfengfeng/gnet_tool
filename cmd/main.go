package main

import (
	"fmt"
	"github.com/Byfengfeng/gnet_tool/net"
	"github.com/Byfengfeng/gnet_tool/utils"
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

	//TestPool()
}

func TestPool()  {

	pool := utils.NewSafePool(func() interface{} {
		return []byte{}
	})
	bytes := pool.Get().([]byte)
	bytes = []byte("123str")
	pool.Put(bytes)
	str := pool.Get().([]byte)
	fmt.Println(string(str))

}