package session

import (
	"database/sql"
	"github.com/futhergo/gee-orm/dialect"
	"github.com/futhergo/gee-orm/log"
	"github.com/futhergo/gee-orm/schema"
	"strings"
)

type Session struct {
	db *sql.DB
	dialect dialect.Dialect
	refTable *schema.Schema
	sql strings.Builder
	sqlVars []interface{}
}

func New(db *sql.DB, d dialect.Dialect) *Session {
	return &Session{
		db: db,
		dialect: d,
		sqlVars: make([]interface{}, 0),
	}
}

func (s *Session)clear() {
	s.sql.Reset()
	s.sqlVars = make([]interface{}, 0)
}

func (s *Session)DB() *sql.DB {
	return s.db
}

func (s *Session)Raw(sql string, vars...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, vars...)
	return s
}

func (s *Session)Exec() (r sql.Result, err error) {
	defer s.clear()
	q := s.sql.String()
	log.Info(q, s.sqlVars)
	r, err = s.db.Exec(q, s.sqlVars...)
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *Session)Query() *sql.Row {
	defer s.clear()
	q := s.sql.String()
	log.Info(q, s.sqlVars)
	return s.db.QueryRow(q, s.sqlVars...)
}

func (s *Session)QueryRows() (rs *sql.Rows, err error) {
	defer s.clear()
	q := s.sql.String()
	rs, err = s.db.Query(q, s.sqlVars...)
	if err != nil {
		log.Error(q, s.sqlVars)
	}
	return
}
