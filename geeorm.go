/*
Engine类

Engine--->DB: 打开/关闭DB，创建Session
    Sessions--->Schema: 解析Model，创建/删除表，执行语句，查询Row/Rows
Dialect：定义不同数据库之间统一的ORM接口，抽离相同的实现，不同的实现由不同类型数据库实现
 */
package gee_orm

import (
	"database/sql"
	"fmt"
	"github.com/futhergo/gee-orm/dialect"
	"github.com/futhergo/gee-orm/log"
	"github.com/futhergo/gee-orm/session"
)

type Engine struct {
	db *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	d, ok := dialect.GetDialect(driver)
	if !ok {
		log.Error( fmt.Sprintf("get %s dialect failed", driver))
		return
	}
	e = &Engine{
		db: db,
		dialect: d,
	}
	log.Info("connect db succeed!")
	return
}

func (e *Engine)Close() {
	if err := e.db.Close(); err != nil {
		log.Error(err)
	}
	log.Info("close db succeed!")
}

func (e *Engine)NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}