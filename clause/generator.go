package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})
type Type int

const (
	INSERT = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderby
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

func _insert(values ...interface{}) (string, []interface{}) {
	// INSERT INTO <tableName> (<fields(columns)>)
	tn := values[0]
	fields := strings.Join(values[1].([]string), ", ")
	return fmt.Sprintf("INSERT INTO %s (%v)", tn, fields), []interface{}{}
}

func generateBindVars(n int) string {
	vars := make([]string, 0)
	for i := 0; i<n; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func _values(values ...interface{}) (string, []interface{}) {
	// "VALUES (?, ?), (?, ?)", "field1V1, field2V1, field1V2, field2V2"
	var bindVars string
	var sql strings.Builder
	var sqlVars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindVars == "" {
			bindVars = generateBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindVars))
		if i != len(values)-1 {
			sql.WriteString(", ")
		}
		sqlVars = append(sqlVars, v...)
	}
	return sql.String(), sqlVars
}

func _select(values ...interface{}) (string, []interface{}) {
	// SELECT <field> FROM <tableName>
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ", ")
	return fmt.Sprintf("SELECT %s FROM %s", fields, tableName), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	// "LIMIT ?", "<num>"
	return "LIMIT ?", values
}

func _where(values ...interface{}) (string, []interface{}) {
	// "WHERE name=?", "Tom"
	return fmt.Sprintf("WHERE %s", values[0]), values[1:]
}

func _orderby(values ...interface{}) (string, []interface{}) {
	// ORDERBY age
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

func _update(values ...interface{}) (string, []interface{}) {
	// UPDATE <tableName> SET field1=?, field2=?, ...
	tableName := values[0]
	kvs := values[1].(map[string]interface{})
	var ks []string
	var vs []interface{}
	for k, v := range kvs {
		ks = append(ks, k+"=?")
		vs = append(vs, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(ks, ", ")), vs
}

func _delete(values ...interface{}) (string, []interface{})  {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"count(*)"})
}
