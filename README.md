[Gee-Orm](https://geektutu.com/post/geeorm.html)学习

# Gee-Orm
数据库ORM库：对数据库的面向对象的抽象
* 创建删除表
* 

# 大纲

* Engine--->DB: 打开/关闭DB，创建Session
* Session--->Schema: 解析Model，创建/删除表，执行语句，查询Row/Rows
* Dialect：定义不同数据库之间统一的ORM接口，抽离相同的实现，不同的实现由不同类型数据库实现