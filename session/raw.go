package session

import (
	"database/sql"
	"github.com/futhergo/gee-orm/clause"
	"github.com/futhergo/gee-orm/dialect"
	"github.com/futhergo/gee-orm/log"
	"github.com/futhergo/gee-orm/schema"
	"strings"
)

type Session struct {
	db *sql.DB
	tx *sql.Tx
	dialect dialect.Dialect
	refTable *schema.Schema
	clause clause.Clause
	sql strings.Builder
	sqlVars []interface{}
}

type CommonDB interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
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
	s.clause = clause.Clause{}
}

func (s *Session)DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
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
	r, err = s.DB().Exec(q, s.sqlVars...)
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *Session)Query() *sql.Row {
	defer s.clear()
	q := s.sql.String()
	log.Info(q, s.sqlVars)
	return s.DB().QueryRow(q, s.sqlVars...)
}

func (s *Session)QueryRows() (rs *sql.Rows, err error) {
	defer s.clear()
	q := s.sql.String()
	rs, err = s.DB().Query(q, s.sqlVars...)
	if err != nil {
		log.Error(q, s.sqlVars)
	}
	return
}
