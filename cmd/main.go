package main

import (
	"fmt"
	"time"
)

func main() {
	//tcpServer := net.NewTcpServer("tcp6","", 9000, true)
	//err := tcpServer.Start()
	//if err != nil {
	//	panic(err)
	//}
	strs := make(chan int)
	go func() {
		for  {
			a := <-strs
			if a == 0 {
				break
			}
			fmt.Println(a)
		}


	}()

	for i:= 0; i < 5; i++ {
		strs <- i
		time.Sleep(time.Second)
	}
	<-strs
	fmt.Println("close chan")
	close(strs)
	time.Sleep(10 * time.Second)
	fmt.Println("close chan")

	<-make(chan struct{})
}