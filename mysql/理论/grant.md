#### grant

##### 创建一个用户
```sql
create user 'kare'@'%' identified by 'kare123';
```
表示创建了一个 `kare` 的用户，密码是 `kare123`.

创建完的用户在 `mysql.user` 表里面.

```sql
select * from mysql.user where user='kare' \G
```

##### 全局权限
```sql
grant all privileges on *.* to 'kare'@'%' with grant option;
```

##### 回收权限
```sql
revoke all privileges on *.* from 'kare'@'%';
```

##### db 权限
支持库级别的定义
```sql
grant all privileges on test1.* to 'ua'@'%' with grant option;
```
指定`db1`的所有权限给用户 `ua`.

基于`库`的权限记录保存在 `mysql.db` 表中.
```sql
select * from mysql.db where user ='kare' \G
```


##### 表权限和列权限

`表权限`定义存放在表 `mysql.tables_priv` 中.`列权限`定义存放在表 `mysql.columns_priv` 中。这两类权限，组合起来存放在内存的 hash 结构 column_priv_hash 中。
```sql
create table test1.tbl_test(id int, a int);
grant all privileges on test1.tbl_test to 'ua'@'%' with grant option;
GRANT SELECT(id), INSERT (id,a) ON mydb.mytbl TO 'ua'@'%' with grant option;
```


`flush privileges` 命令会清空 `acl_users` 数组，然后从 mysql.user` 表中读取数据`重新加载`，重新构造一个 acl_users 数组。也就是说，以数据表中的数据为准，会将全局权限内存数组重新加载一遍


1. `grant` 语句会同时修改`数据表`和`内存`，`判断权限的时候使用的是内存数据`。因此，规范地使用 grant 和 revoke 语句，是不需要随后加上 flush privileges 语句的。

2. flush privileges` 语句本身会`用数据表的数据重建一份内存权限数据，所以在`权限数据可能存在不一致的情况下再使用`。而这种不一致往往是由于直接用 DML 语句操作系统权限表导致的，所以我们尽量不要使用这类语句


* 小知识点
1. \G 行转列并发送给 mysql server
2. \s mysql status 信息
3. \h 显示可用帮助



