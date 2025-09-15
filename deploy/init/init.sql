-- 创建 root@'%' 用户（如已存在则跳过）
CREATE USER IF NOT EXISTS 'root'@'%' IDENTIFIED BY '2214380963Wx!!';

-- 授权所有权限
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION;

-- 应用权限修改
FLUSH PRIVILEGES;

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `coderhub`;