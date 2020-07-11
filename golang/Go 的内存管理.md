众所周知，Go 是一门自带垃圾回收机制的语言，内存管理参照 tcmalloc 实现，使用连续虚拟地址，以页( 8k )为单位、多级缓存进行管理。针对小于16 byte 直接使用Go的上下文P中的mcache分配，大于 32 kb 直接在 mheap 申请，剩下的先使用当前 P 的 mcache 中对应的 size class 分配 ，如果 mcache 对应的 size class 的 span 已经没有可用的块，则向 mcentral 请求。如果 mcentral 也没有可用的块，则向 mheap 申请，并切分。如果 mheap 也没有合适的 span，则向操作系统申请。


Go 在内存统计方面做的也是相当出色，提供细粒度的内存分配、GC 回收、goroutine 管理等统计数据


go_memstats_sys_bytes ：进程从操作系统获得的内存的总字节数 ，其中包含 Go 运行时的堆、栈和其他内部数据结构保留的虚拟地址空间。  
go_memstats_heap_inuse_bytes：在 spans 中正在使用的字节。其中不包含可能已经返回到操作系统，或者可以重用进行堆分配，或者可以将作为堆栈内存重用的字节。
go_memstats_heap_idle_bytes：在 spans 中空闲的字节。
go_memstats_stack_sys_bytes：栈内存字节，主要用于 goroutine 栈内存的分配



Go 提供了三种内存回收机制：定时触发，按量触发，手动触发