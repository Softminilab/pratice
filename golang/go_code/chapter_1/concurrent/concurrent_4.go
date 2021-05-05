package main

//
//import (
//	"fmt"
//	"time"
//)
//
//func main() {
//	in := make(chan int)
//	out := make(chan int)
//	go producer(time.Millisecond*100, in)
//	go producer(time.Millisecond*250, in)
//	go consumer(out)
//	for i := range in {
//		out <- i
//	}
//}
//
//func producer(t time.Duration, in chan int) {
//	var i int
//	for {
//		in <- i
//		i++
//		time.Sleep(t)
//	}
//}
//
//func consumer(out chan int) {
//	for t := range out {
//		fmt.Println(t)
//	}
//}
