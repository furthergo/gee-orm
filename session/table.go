package session

import (
	"fmt"
	"github.com/futhergo/gee-orm/log"
	"github.com/futhergo/gee-orm/schema"
	"reflect"
	"strings"
)

func (s *Session)Model(v interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(s.refTable.Model) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(v, s.dialect)
	}
	return s
}

func (s *Session)RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("session refTable is empty")
	}
	return s.refTable
}

func (s *Session)CreateTable() error {
	t := s.RefTable()
	var cols []string
	for _, f := range t.Fields {
		cols = append(cols, fmt.Sprintf("%s %s %s", f.Name, f.Type, f.Tag))
	}
	desc := strings.Join(cols, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", t.Name, desc)).Exec()
	return err
}

func (s *Session)DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

func (s *Session)HasTable() bool {
	q, args := s.dialect.TableExistSQL(s.RefTable().Name)
	r := s.Raw(q, args...).Query()
	var t string
	_ = r.Scan(&t)
	return t == s.RefTable().Name
}