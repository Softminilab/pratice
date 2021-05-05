package main

import "net"

//func main() {
//	var wg sync.WaitGroup
//	wg.Add(36)
//	pool(&wg, 36, 50)
//	wg.Wait()
//}
//
//func worker(tasks <-chan int, wg *sync.WaitGroup) {
//	defer wg.Done()
//
//	for {
//		task, ok := <-tasks
//		if !ok {
//			return
//		}
//		d := time.Duration(task) * time.Millisecond
//		time.Sleep(d)
//		fmt.Println("processing task ", task)
//	}
//}
//
//func pool(wg *sync.WaitGroup, workers, tasks int) {
//	taskCh := make(chan int)
//	for i := 0; i < workers; i++ {
//		go worker(taskCh, wg)
//	}
//
//	for i := 0; i < tasks; i++ {
//		taskCh <- i
//	}
//
//	close(taskCh)
//}

func handler(c net.Conn) {
	c.Write([]byte("ok"))
	c.Close()
}

func main() {
	l, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}
		go handler(c)
	}
}
