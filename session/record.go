package session

import (
	"errors"
	"github.com/futhergo/gee-orm/clause"
	"reflect"
)

func (s *Session)Insert(values ...interface{}) (int64, error) {
	if len(values) == 0 {
		return 0, nil
	}
	vs := make([]interface{}, 0)
	schema := s.Model(values[0]).RefTable()
	s.clause.Set(clause.INSERT, schema.Name, schema.FieldsName)
	for _, value := range values {
		s.CallMethod(BeforeInsert, value)
		vs = append(vs, schema.RecordValues(value))
	}
	s.clause.Set(clause.VALUES, vs...)
	q, v := s.clause.Build(clause.INSERT, clause.VALUES)
	res, err := s.Raw(q, v...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterInsert, nil)
	return res.RowsAffected()
}

func (s *Session)Find(values interface{}) error {
	s.CallMethod(BeforeQuery, nil)
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType :=destSlice.Type().Elem()
	m := reflect.New(destType).Elem().Interface()
	table := s.Model(m).RefTable()
	s.clause.Set(clause.SELECT, table.Name, table.FieldsName)
	q, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(q, vars...).QueryRows()
	if err != nil {
		return err
	}
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var vs []interface{}
		for _, name := range table.FieldsName {
			vs = append(vs, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(vs...); err != nil {
			return err
		}
		s.CallMethod(AfterQuery, dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

func (s *Session)Update(kvs ...interface{}) (int64, error) {
	s.CallMethod(BeforeUpdate, nil)
	m, ok := kvs[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i<len(kvs); i+=2 {
			m[kvs[i].(string)] = kvs[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	q, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	res, err := s.Raw(q, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterUpdate, nil)
	return res.RowsAffected()
}

func (s *Session)Delete() (int64, error) {
	s.CallMethod(BeforeDelete, nil)
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	q, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	res, err := s.Raw(q, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterDelete, nil)
	return res.RowsAffected()
}

func (s *Session)Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	q, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(q, vars...).Query()

	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Session)Where(desc string, args ...interface{}) *Session {
	// WHERE field=?, fieldV
	s.clause.Set(clause.WHERE, append([]interface{}{desc}, args...)...)
	return s
}

func (s *Session)Limit(n int64) *Session {
	s.clause.Set(clause.LIMIT, n)
	return s
}

func (s *Session)OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

func (s *Session)First(m interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(m))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("Record Not Found")
	}
	dest.Set(destSlice.Index(0))
	return nil
}