#### join 中的几种算法

创建两个表
```sql
CREATE TABLE `t2` (
  `id` int(11) NOT NULL,
  `a` int(11) DEFAULT NULL,
  `b` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `a` (`a`)
) ENGINE=InnoDB;

drop procedure idata;
delimiter ;;
create procedure idata()
begin
  declare i int;
  set i=1;
  while(i<=1000)do
    insert into t2 values(i, i, i);
    set i=i+1;
  end while;
end;;
delimiter ;
call idata();

create table t1 like t2;
insert into t1 (select * from t2 where id<=100)
```

##### NLJ【Index Nested-Loop Join】
执行语句
```sql
select * from t1 straight_join t2 on (t1.a=t2.a);
```
`straight_join` 指定 t1 为驱动表，t2 为被驱动表。让优化器按照我们指定的方式去join

###### 在用上索引的情况下，join 的流程
1. 从表 t1 中取出一行 R；
2. 从数据 R 中取出 a 字段，去 t2 中去查找；
3. 取出 t2中满足条件的行，和 R 组成一行，作为结果集的一部分
4. 重复 1 到 3 步骤，直到循环到 t1 表的末尾为止；

这个过程是先遍历表 t1，然后根据从表 t1 中取出的每行数据中的 a 值，去表 t2 中查找满足条件的记录。在形式上，这个过程就跟我们写程序时的嵌套查询类似，并且可以用上`被驱动表的索引`，所以我们称之为“Index Nested-Loop Join”，简称 NLJ。

##### Simple Nested-Loop Join
如果 sql 改成这样
```sql
select * from t1 straight_join t2 on (t1.a=t2.b);
```
由于表 t2 的字段 b 上没有索引，每次到 t2 去匹配的时候，就要做一次`全表扫描`.这算法名称叫做`Simple Nested-Loop Join` ,mysql 并没有使用这个算法，而是使用了一个叫`Block Nested-Loop Join` 的算法，简称 `BNL`.

##### Block Nested-Loop Join【BNL】
这时候，被驱动表上`没有可用的索引`，算法的流程是这样的：

1. 把表 `t1` 的数据读入`线程内存 join_buffer` 中，由于我们这个语句中写的是 select *，因此是把整个表 `t1 放入了内存`；
2. 扫描表 t2，把表 t2 中的`每一行取出来`，跟 `join_buffer 中的数据做对比`，满足 join 条件的，作为结果集的一部分返回

`join_buffer_size` 默认是 256k,如果是一个大表，数据很多，那么一次是不能全部读入内存的。这个时候就要`分段`放入了.
```sql
select * from t1 straight_join t2 on (t1.a=t2.b);
```
执行过程就是：

1. 扫描表 t1，顺序读取数据行放入 join_buffer 中，放完第 88 行 join_buffer 满了，继续第 2 步；
2. 扫描表 t2，把 t2 中的每一行取出来，跟 join_buffer 中的数据做对比，满足 join 条件的，作为结果集的一部分返回；
3. 清空 join_buffer；
4. 继续扫描表 t1，顺序读取最后的 12 行数据放入 join_buffer 中，继续执行第 2 步。

![Block Nested-loop join 流程图](https://github.com/karepbq/pratice/blob/master/mysql/%E7%90%86%E8%AE%BA/img/Block_Nested_Loop_Join.jpg)

图中的步骤 4 和 5，表示清空 join_buffer 再复用。
这个流程才体现出了这个算法名字中“Block”的由来，表示“分块去 join”

`如果你的 join 语句很慢，就把 join_buffer_size 改大。`

* 第一个问题：能不能使用 join 语句？
1. 如果可以使用 Index Nested-Loop Join 算法，也就是说可以用上被驱动表上的索引，其实是没问题的；
2. 如果使用 Block Nested-Loop Join 算法，扫描行数就会过多。尤其是在大表上的 join 操作，这样可能要扫描被驱动表很多次，会占用大量的系统资源。所以这种 join 尽量不要用

##### 驱动表改如何选择
在决定哪个表做驱动表的时候，应该是两个表按照各自的`条件过滤`，过滤完成之后，计算参与 join 的各个`字段的总数据量`，`数据量小的那个表`，`就是“小表”，应该作为驱动表`

--------

#### join 的优化

```sql
create table t1(id int primary key, a int, b int, index(a));
create table t2 like t1;
drop procedure idata;
delimiter ;;
create procedure idata()
begin
  declare i int;
  set i=1;
  while(i<=1000)do
    insert into t1 values(i, 1001-i, i);
    set i=i+1;
  end while;
  
  set i=1;
  while(i<=1000000)do
    insert into t2 values(i, i, i);
    set i=i+1;
  end while;

end;;
delimiter ;
call idata();
```

##### Multi-Range Read 优化【MRR】
MRR的主要目的是`尽量按照顺序读盘`，`检查随机访问`。

MRR 的设计思想：`大多数的数据都是按照主键递增顺序插入的得到的，所以我们可以认为，如果按照主键的递增顺序查询的话，对磁盘的读比较接近顺序读，能够提升性能`

如果不是按照 MRR 的方法查询一条语句是怎么样的了？
```sql
select * from t1 where a>=1 and a<=100;
```

例如上面的语句，我们就需要根据 a 的值，找到 a 所在位置的主键 id 值后，去主键索引上查询真实的值，然后返回。也就是回表过程

如果使用 MRR 了，查询流程就变成这样了
1. 根据索引 a 定位满足条件的记录，将 id 值存入 read_rnd_buffer 中
2. 将 `read_rnd_buffer` 中的值根据 id 就行递增排序
3. 排序后的 id 数组，依次到主键 id 索引中查记录，并作为结果返回

这里，`read_rnd_buffer` 的大小是由 `read_rnd_buffer_size` 参数控制的。如果步骤 1 中，`read_rnd_buffer 放满了`，就会先执行完步骤 `2` 和 `3`，然后清空 `read_rnd_buffer`。之后继续找索引 a 的下个记录，并继续循环

你想要稳定地使用 MRR 优化的话，需要设置`set optimizer_switch="mrr_cost_based=off"`。（官方文档的说法，是现在的优化器策略，判断消耗的时候，会更倾向于不使用 MRR，把 mrr_cost_based 设置为 off，就是固定使用 MRR 了。）

##### Batched Key Access 【BKA】
其实 BKA 算法是会对 NLJ 算法的优化。

NLJ的执行逻辑是：从驱动表 t1，`一行行`地取出 a 的值，再到被驱动表 t2 去做 join。也就是说，对于表 t2 来说，每次`都是匹配一个值`

那么 BKA 的思想就是一次从驱动表 t1，`取多行值`，放入 `join_buffer`，简单的说就是从之前的取一行变成一次取多行，放入 join_buffer.让被驱动表走索引。

开启 BKA
```sql
set optimizer_switch='mrr=on,mrr_cost_based=off,batched_key_access=on';
```
前两个参数的作用是要启用 `MRR`。这么做的原因是，`BKA` 算法的优化要`依赖`于 `MRR`



