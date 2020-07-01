同级目录的 my.cnf

装好了mysql登录方式为：
`mysql -u root -h 127.0.0.1 -p`

然后输入系统随机生成的那密码


#### mysql初始化
`./bin/mysqld --user=root --basedir=/usr/local/mysql/mysql-5.7.30 --datadir=/usr/local/mysql/mysql-5.7.30/data --initialize`

#### 删除mysql
https://gist.github.com/vitorbritto/0555879fe4414d18569d

#### 如何手动安装
https://www.cnblogs.com/jonney-wang/p/11279220.html


https://gist.github.com/zubaer-ahammed/c81c9a0e37adc1cb9a6cdc61c4190f52

#### 更改密码
1. sudo mysqld_safe --skip-grant-tables
2. mysql -u root
3. UPDATE mysql.user SET authentication_string=null WHERE User='root';
4. FLUSH PRIVILEGES;
5. exit

then

`mysql -u root`

`ALTER USER 'root'@'localhost' IDENTIFIED WITH caching_sha2_password BY 'yourpasswd'`;
or
`ALTER USER 'root'@'localhost' IDENTIFIED BY '123456'`;


#### 安全停止mysql
`/usr/local/mysql/bin/mysqladmin -u root -p shutdown`