package session

import (
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
		vs = append(vs, schema.RecordValues(value))
	}
	s.clause.Set(clause.VALUES, vs...)
	q, v := s.clause.Build(clause.INSERT, clause.VALUES)
	res, err := s.Raw(q, v...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *Session)Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType :=destSlice.Type().Elem()
	m := reflect.New(destType).Elem().Interface()
	table := s.Model(m).RefTable()
	s.clause.Set(clause.SELECT, table.Name, table.FieldsName)
	q, vars := s.clause.Build(clause.SELECT)
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
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}