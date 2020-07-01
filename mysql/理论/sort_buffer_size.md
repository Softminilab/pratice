`sort_buffer_size`，就是 MySQL 为`排序`开辟的`内存（sort_buffer）的大小`。如果要排序的数据量小于 sort_buffer_size，排序就在`内存中完成`。但如果`排序数据量太大`，内存放不下，`则不得不利用磁盘临时文件辅助排序`

```sql
CREATE TABLE `t` (
  `id` int(11) NOT NULL,
  `city` varchar(16) NOT NULL,
  `name` varchar(16) NOT NULL,
  `age` int(11) NOT NULL,
  `addr` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `city` (`city`)
) ENGINE=InnoDB;

insert into t value (1,'杭州','张三',123,'杭州');
```

```sql
/* 打开optimizer_trace，只对本线程有效 */
SET optimizer_trace='enabled=on'; 

/* @a保存Innodb_rows_read的初始值 */
select VARIABLE_VALUE into @a from  performance_schema.session_status where variable_name = 'Innodb_rows_read';

/* 执行语句 */
select city, name,age from t where city='杭州' order by name limit 1000; 

/* 查看 OPTIMIZER_TRACE 输出 */
SELECT * FROM `information_schema`.`OPTIMIZER_TRACE`\G

/* @b保存Innodb_rows_read的当前值 */
select VARIABLE_VALUE into @b from performance_schema.session_status where variable_name = 'Innodb_rows_read';

/* 计算Innodb_rows_read差值 */
select @b-@a;
```

这个方法是通过查看 OPTIMIZER_TRACE 的结果来确认的，你可以从 number_of_tmp_files 中看到是否使用了临时文件

![DML](https://github.com/karepbq/pratice/blob/master/mysql/%E7%90%86%E8%AE%BA/img/optimizer_trace.png)


`number_of_tmp_files` 表示的是，`排序过程中使用的临时文件数`。你一定奇怪，为什么需要 12 个文件？内存放不下时，就需要使用`外部排序`，外部排序一般使用`归并排序算法`。可以这么简单理解，MySQL 将需要排序的数据分成 12 份，每一份单独排序后存在这些临时文件中。然后把这 12 个有序文件再合并成一个有序的大文件。

如果 `sort_buffer_size` `超过了需要排序的数据量的大小`，`number_of_tmp_files` 就是 `0`，表示`排序可以直接在内存中完成。`

否则就需要放在临时文件中排序。`sort_buffer_size 越小`，`需要分成的份数越多`，n`umber_of_tmp_files 的值就越大`


`max_length_for_sort_data`，是 MySQL 中专门控制用于排序的行数据的长度的一个参数。它的意思是，如果单行的长度超过这个值，MySQL 就认为单行太大，要换一个算法