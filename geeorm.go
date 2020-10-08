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
	"strings"
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

func (e *Engine)Transaction(fn func(s *session.Session) (interface{}, error)) (v interface{}, err error) {
	s := e.NewSession()
	err = s.Begin()
	if err != nil {
		return nil, err
	}
	paniced := true
	defer func() {
		// 参考gorm的做法，fn出错、commit出错或者panic时，rollback事务
		if paniced || err != nil {
			s.Rollback()
		}
	}()

	v, err = fn(s)
	if err == nil {
		err = s.Commit()
	}
	paniced = false
	return v, err
}

func difference(target, origin []string) []string {
	// target - origin
	res := make([]string, 0)
	m := make(map[string]bool)
	for _, v := range origin {
		m[v] = true
	}

	for _, v := range target {
		if _, ok := m[v]; !ok {
			res = append(res, v)
		}
	}
	return res
}

func (e *Engine)Migrate(value interface{}) error {
	_, err := e.Transaction(func(s *session.Session) (interface{}, error) {
		schema := s.Model(value).RefTable()
		if !s.HasTable() {
			log.Info(fmt.Sprintf("Migrate of not exist table %s", schema.Name))
			s.CreateTable()
		}
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", schema.Name)).QueryRows()
		columns, _ := rows.Columns()
		addColumns := difference(schema.FieldsName, columns)
		delColumns := difference(columns, schema.FieldsName)

		log.Info(fmt.Sprintf("add columns %v, delete columns %s", addColumns, delColumns))

		for _, c := range addColumns {
			_, err := s.Raw(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", schema.Name, c, schema.FieldsMap[c].Type)).Exec()
			if err != nil {
				return nil, err
			}
		}

		if len(delColumns) == 0 {
			return nil, nil
		}

		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s FROM %s;", "deleteTable_" + schema.Name, strings.Join(schema.FieldsName, ", "), schema.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", schema.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", "deleteTable_" + schema.Name, schema.Name))
		_, err := s.Exec()
		return nil, err
	})
	return err
}
