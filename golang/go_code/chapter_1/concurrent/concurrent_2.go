package main
//
//import (
//	"fmt"
//	"time"
//)
//
//func main() {
//	for i:=0;i<24;i++ {
//		c := timer(time.Second * 1)
//		fmt.Println(<-c)
//	}
//}
//
//func timer(duration time.Duration) <-chan int {
//	c := make(chan int)
//	go func() {
//		time.Sleep(duration)
//		c <- 1
//	}()
//	return c
//}