package main

import (
	"fmt"
	"time"
)

func main() {
	//var wg sync.WaitGroup
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	time.Sleep(time.Second * 5)
	//	fmt.Println("时间到了")
	//}()
	//fmt.Println("111")
	//wg.Wait()
	//fmt.Println("2222")
	//done := make(chan bool)
	//
	//values := []string{"a", "b", "c"}
	//for _, v := range values {
	//	go func(x string) {
	//		fmt.Println(x)
	//		done <- true
	//	}(v)
	//}
	//
	//// wait for all goroutines to complete before exiting
	//for _ = range values {
	//	fmt.Println(<-done)
	//}

	ch := make(chan int)
	go func() {
		time.Sleep(time.Second * 1)
		fmt.Println("sort data 2")
		ch <- 1
	}()
	fmt.Println("sort data 1")
	<-ch
	fmt.Println("merging")
}
