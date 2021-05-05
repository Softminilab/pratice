package main

//
//import (
//	"fmt"
//	"time"
//)
//
//func main() {
//	var Ball int
//	table := make(chan int)
//	for i := 0; i < 100; i++ {
//		fmt.Printf("%v %p \n", i, &i)
//		go player(table)
//	}
//
//	table <- Ball
//	time.Sleep(1 * time.Second)
//	fmt.Println(<-table)
//}
//
//func player(table chan int) {
//	for {
//		ball := <-table
//		ball++
//		time.Sleep(100 * time.Millisecond)
//		table <- ball
//	}
//}
