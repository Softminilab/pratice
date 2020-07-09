### redo log & binlog
他们是两种不同的日志 log，两者的区别如下
1. `redo log` 是 `InnoDB` 引擎持有的；`binlog` 是 `mysql` 的 `server` 层实现的，所有引擎都可以使用.
2. `redo log` 是`物理日志`，记录的是“某个数据页上做了什么修改”；`binlog` 是`逻辑日志`，记录的是这个语句的原始逻辑，比如“给 ID2 的这一样的 C 字段加 1”。
3. `redo log` 是`循环写`的，`空间固定会使用完`；`binlog` 是可以`追加写入`的，“追加写”是指 binlog 文件写到一定大小后会切换下一个，`并不会覆盖以前的日志`。

#### binlog 的写入机制
binlog 的写入逻辑比较简单：`事务执行过程中，先把日志写到 binlog cache，事务提交的时候，再把 binlog cache 写到 binlog 文件中`。

系统给 `binlog cache` 分配了一片`内存`，每个线程一个，参数 `binlog_cache_size` 用于控制`单个线程`内 `binlog cache` 所占内存的大小。如果超过了这个参数规定的大小，就要暂存到`磁盘`。

![binlog写盘状态](https://github.com/kareTauren/pratice/blob/master/mysql/%E7%90%86%E8%AE%BA/img/binlog.png)

每个线程有`自己 binlog cache`，但是共用同一份 `binlog` 文件。

* 图中的 write，指的就是指把日志写入到文件系统的 page cache，并没有把数据持久化到磁盘，所以速度比较快。

* 图中的 fsync，才是将数据持久化到磁盘的操作。一般情况下，我们认为 fsync 才占磁盘的 IOPS。

write 和 fsync 的时机，是由参数 sync_binlog 控制的
1. sync_binlog=0 的时候，表示每次提交事务都只 write，不 fsync；
2. sync_binlog=1 的时候，表示每次提交事务都会执行 fsync；
3. sync_binlog=N(N>1) 的时候，表示每次提交事务都 write，但累积 N 个事务后才 fsync。

因此，在出现 IO 瓶颈的场景里，将 `sync_binlog` 设置成一个比较`大`的值，`可以提升性能`。在实际的业务场景中，考虑到丢失日志量的可控性，一般不建议将这个参数设成 0，比较常见的是将其设置为 `100~1000` 中的某个数值。

但是，将 sync_binlog 设置为 N，对应的风险是：`如果主机发生异常重启，会丢失最近 N 个事务的 binlog 日志。`

`binlog 是不能“被打断的”。一个事务的 binlog 必须连续写，因此要整个事务完成后，再一起写到文件里`

#### redo log 的写入机制

![redo log 状态](https://github.com/kareTauren/pratice/blob/master/mysql/%E7%90%86%E8%AE%BA/img/redolog_status.png)

这三种状态分别是：
1. 存在 `redo log buffer` 中，物理上是在 MySQL `进程内存中`，就是图中的红色部分；
2. 写到磁盘 `(write)`，但是`没有持久化（fsync)`，物理上是在`文件系统的` `page cache` 里面，也就是图中的黄色部分；
3. `持久化到磁盘`，对应的是 `hard disk`，也就是图中的绿色部分

日志写到 redo log buffer 是很快的，wirte 到 page cache 也差不多，但是持久化到磁盘的速度就慢多了

为了控制 redo log 的写入策略，InnoDB 提供了 innodb_flush_log_at_trx_commit 参数，它有三种可能取值：
1. 设置为 `0` 的时候，表示每次事务提交时都只是把 `redo log` 留在 `redo log buffer` 中 ;
2. 设置为 `1` 的时候，表示每次事务提交时都将 `redo log 直接持久化到磁盘`；
3. 设置为 `2` 的时候，表示每次事务提交时都只是把 `redo log 写到 page cache`。

两阶段提交的时候说过，时序上 `redo log` 先 `prepare`， 再写 `binlog`，最后再把 `redo log commit`


#### binlog 的三种格式
1. `statement` 只是记录大概的操作过程，
2. `row` 会记录详细的信息到 binlog 文件
3. `mixed` 是介于 row 和 staement 两种格式之间的一种格式

为什么会有 mixed 格式的 binlog？

* 因为有些 statement 格式的 binlog 可能会导致`主备不一致`，所以要使用 `row` 格式。

* 但 row 格式的缺点是，`很占空间`。比如你用一个 delete 语句删掉 10 万行数据，用 statement 的话就是一个 SQL 语句被记录到 binlog 中，占用几十个字节的空间。但如果用 row 格式的 binlog，就要把这 10 万条记录都写到 binlog 中。这样做，不仅会占用更大的空间，`同时写 binlog 也要耗费 IO 资源，影响执行速度`。

* 所以，MySQL 就取了个折中方案，也就是有了 mixed 格式的 binlog。mixed 格式的意思是，MySQL 自己会判断这条 SQL 语句是否可能引起主备不一致，如果有可能，就用 row 格式，否则就用 statement 格式

也就是说，mixed 格式可以利用 statment 格式的优点，同时又避免了数据不一致的风险。



