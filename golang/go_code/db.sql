-- create db
create database go;

-- creat table
CREATE TABLE `tbl_users` (
    `id` bigint(20) NOT NULL COMMENT 'user id',
    `name` varchar(32) NOT NULL DEFAULT '' COMMENT 'user name',
    `email` varchar(32) NOT NULL DEFAULT '' COMMENT 'user email',
    PRIMARY KEY (`id`),
    UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 comment 'user table';