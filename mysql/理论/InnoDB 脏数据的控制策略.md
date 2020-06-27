#### InnoDB 脏页的控制策略，以及和这些策略相关的参数

脏页: 当内存数据页跟磁盘数据页内容不一致的时候，我们称这个内存页为。
干净页: 内存数据写入到磁盘后，内存和磁盘上的数据页的内容就一致了。


这就要用到 innodb_io_capacity 这个参数了，它会告诉 InnoDB 你的磁盘能力。这个值我建议你设置成磁盘的 IOPS。磁盘的 IOPS 可以通过 fio 这个工具来测试，下面的语句是我用来测试磁盘随机读写的命令
```
 fio -filename=$filename -direct=1 -iodepth 1 -thread -rw=randrw -ioengine=psync -bs=16k -size=500M -numjobs=10 -runtime=10 -group_reporting -name=mytest  
```

其实，因为没能正确地设置 innodb_io_capacity 参数，而导致的性能问题.

参数 `innodb_max_dirty_pages_pct` 是脏页比例上限，默认值是 `75%`

合理地设置 innodb_io_capacity 的值，并且平时要多`关注脏页比例，不要让它经常接近 75%`。

其中，脏页比例是通过 
`Innodb_buffer_pool_pages_dirty/Innodb_buffer_pool_pages_total` 得到的，具体的命令参考下面的代码：
```sql
mysql> select VARIABLE_VALUE into @a from global_status where VARIABLE_NAME = 'Innodb_buffer_pool_pages_dirty';
select VARIABLE_VALUE into @b from global_status where VARIABLE_NAME = 'Innodb_buffer_pool_pages_total';
select @a/@b;
```