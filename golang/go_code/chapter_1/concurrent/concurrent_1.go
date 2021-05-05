package main
//
//import (
//	"bytes"
//	"fmt"
//	"runtime"
//	"strconv"
//)
//
//func main() {
//	// 创建一个int类型的通道
//	c := make(chan int)
//
//	// 开启一个匿名的goroutine
//	go func() {
//		// 向通道发送数字1
//		c <- 1
//	}()
//
//	// 从通道读取结果
//	rs := <-c
//	fmt.Println(rs)
//
//	fmt.Println(getGID())
//}
//
//func getGID() uint64 {
//	b := make([]byte, 64)
//	b = b[:runtime.Stack(b, false)]
//	b = bytes.TrimPrefix(b, []byte("goroutine "))
//	b = b[:bytes.IndexByte(b, ' ')]
//	n, _ := strconv.ParseUint(string(b), 10, 64)
//	return n
//}
//
//
