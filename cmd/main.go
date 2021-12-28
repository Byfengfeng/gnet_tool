package main

import "fmt"

type S struct {
	M *int
}
func main() {
	//var p map[int]int
	//p = make(map[int]int)
	////l := sync.Mutex{}
	//var wg = &sync.WaitGroup{}
	//wg.Add(1)
	//for i := 0; i < 2; i++ {
	//	a := i
	//
	//	go func() {
	//		wg.Wait()
	//		defer wg.Done()
	//		p[a] = 1
	//	}()
	//
	//}
	//
	//for i := 0; i < 2; i++ {
	//	fmt.Println(p[i])
	//}

	var x S
	y := &x
	identity(y)
	//a := new(int)
	//b := t{}
	//structtest(&b)
	//refStruct(a)
	//var a1 int = 1
	//refStruct2(&a1)
	var a *int
	c := newT(a)
	fmt.Println(c)
	trick()
	//*i = 5
	//j := refStruct()
	//fmt.Println(*i)
	//fmt.Println(*j)
	//tcp.NewNetListen(":8999",network.NewNetWork())
	//<-make(chan struct{})

	//TestPool()
}
func identity(z *S) *S {
	return z
}
func newT(a *int) *int {
	return a
}
func refStruct(a *int) *int {
	*a = 1
	return a
}

func refStruct2(a *int) *int {
	return a
}

//func structtest(t *t){
//}

type t struct {
	id int32
}