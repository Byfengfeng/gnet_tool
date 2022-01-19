package main

import (
	"fmt"
	"github.com/Byfengfeng/gnet_tool/utils"
	"github.com/valyala/bytebufferpool"
	"sort"
	"time"
)

type S struct {
	M *int
}

func main() {
	byteBuffer := bytebufferpool.Get()
	//byteBuffer.Write()
	str := "123"
	byteBuffer.Write([]byte(str))
	byteBuffer.Write([]byte("str"))
	fmt.Println("byteBuffer:",byteBuffer.Bytes())
	bytebufferpool.Put(byteBuffer)
	byteBuffer1 := bytebufferpool.Get()
	fmt.Println("byteBuffer1:",byteBuffer1.Bytes())
	//listen := tcp.NewNetListen("192.168.31.134:9000")
	//listen.Start()
	//fmt.Println(splitSort([]int{1,3,5,9,11,65,78,99},11))
	bytes := utils.NewBytes( 1024, func(bytes []byte) {
		//decode, data := utils.Decode(bytes)
		//fmt.Println("decode:",decode,"data:",data)
		fmt.Println(string(bytes))
	})
	//go bytes.ReadBytes()
	go func() {
		for i := 0; i < 1000; i++ {
			time.Sleep(1)
			if i % 2 == 0 {
				str := "00000"
				lens := uint16(len(str))
				bys := make([]byte,0)
				bys = append(bys,byte(lens >> 8),byte(lens))
				bys = append(bys,[]byte(str)...)
				bytes.WriteBytes(uint16(len(bys)),bys)
			}else{
				str := "1111111111"
				lens := uint16(len(str))
				bys := make([]byte,0)
				bys = append(bys,byte(lens >> 8),byte(lens))
				bys = append(bys,[]byte(str)...)
				bytes.WriteBytes(uint16(len(bys)),bys)
			}
		}
		bytes.Len()
		//for i := 0; i < 1000; i++ {
		//	bytes.ReadBytes()
		//}

	}()


	time.Sleep(1 * time.Minute)
	//data := []int{1,  2, 4, 12, 21, 8, 12, 31, 24, 12, 14, 23}
	//listData := QuickSort(data)
	//fmt.Println(listData)
	//source := []int{1, 2, 4, 5, 6}
	//key := 5
	//index := splitSort(source,key)
	//fmt.Println(index)
	//forDel()
	//var arr = []int{2, 3, 4, 1}
	//change(arr)
	//fmt.Println(arr)
	//var p map[int]int
	//if p == nil {
	//
	//}
	//p = make(map[int]int,2)
	//////l := sync.Mutex{}
	//var wg = &sync.WaitGroup{}
	//
	//for i := 0; i < 2; i++ {
	//	//a := i
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		p[i] = 1
	//	}()
	//	wg.Wait()
	//}
	//
	//for i := 0; i < 2; i++ {
	//	fmt.Println(p[i])
	//}

	//var x S
	//y := &x
	//identity(y)
	//a := new(int)
	//b := t{}
	//structtest(&b)
	//refStruct(a)
	//var a1 int = 1
	//refStruct2(&a1)
	//var a *int
	//c := newT(a)
	//fmt.Println(c)、
	//*i = 5
	//j := refStruct()
	//fmt.Println(*i)
	//fmt.Println(*j)
	//tcp.NewNetListen(":8999",network.NewNetWork())
	<-make(chan struct{})

	//TestPool()
}

func change(arr []int)  {
	arr = append(arr,5)
	sort.Ints(arr)
	fmt.Println(arr)

}

func QuickSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	startNum := arr[0]
	left := make([]int,0)
	center := make([]int,0)
	right := make([]int,0)
	for _,i := range arr {
		if i < startNum {
			left = append(left,i)
		}else if i > startNum {
			right = append(right,i)
		}else{
			center = append(center,i)
		}
	}
	left,right = QuickSort(left),QuickSort(right)
	return append(append(left,center...),right...)
}

func forDel()  {
	data := []int{1,  2, 4, 12, 21, 8, 12, 31, 24, 12, 14, 23}
	deleteNum := 2
	for i := 0; i < len(data); i++ {
		if data[i] == deleteNum {
			data = append(data[:i],data[i+1:]...)
		}
	}
}

func splitSort(array []int,target int) int {
	start,end := 0,len(array)-1
	for start <= end {
		index := start + (end - start) / 2
		if array[index] == target {
			return index
		}
		if array[index] > target {
			end = index - 1
		}else{
			start = index + 1
		}
	}
	return -1
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
