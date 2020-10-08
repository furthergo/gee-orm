[Gee-Orm](https://geektutu.com/post/geeorm.html)学习

# Gee-Orm
数据库ORM库：对数据库的面向对象的抽象
* 创建删除表：解析Model，Fields，Tag
* Find/Insert Model

# 大纲

* Engine--->DB: 打开/关闭DB，创建Session
* Session--->Schema: 解析Model，创建/删除表，执行语句，查询Row/Rows，Insert/Find
* Dialect：定义不同数据库之间统一的ORM接口，抽离相同的实现，不同的实现由不同类型数据库实现

重点：
1. reflect
2. package设计

## 功能模块

* Clause生成子句：INSERT/VALUES/SELECT/WHERE/LIMIT/ORDERBY
* session/Record.go 封装ORM操作，
    * Insert：Insert时根据Model，用反射得到Schema和相关字段名和字段值，设置给Clause生成子句，然后执行
            * s.Insert(&u1, &u2) => `"INSERT INTO User (Name text, Age integer) VALUES (?, ?), (?, ?)", ("Tom", 15), ("Jack", 23)`;
            * Model => tableName
            * Model{} => Values: ("Tom", 15), ("Jack", 23)
            * INSERT INTO <tableName> VALUES (?, ?) ... Values.Join(", ");
    * Find：根据传入的切片类型，用反射生成Model，解析Schema，查询Rows，遍历得到真正的值添加回切片
    * Update/Delete/Count
    * 函数返回指针本身链式调用
* session/hooks：封装实现hook，定义特定的函数名，通过反射找到对应的函数并在合适的时机调用
* session/transaction: 封装事务操作
* geeorm.go：封装Engine相关的操作，一个Engine对应一个DB连接，生成N个Session用于操作