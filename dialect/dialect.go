package dialect

import (
	"reflect"
	"strings"
)

var dialectMaps = make(map[string]Dialect)

type Dialect interface {
	DataTypeOf(tpe reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{})
}

func RegisterDialect(name string, d Dialect) {
	dialectMaps[strings.ToLower(name)] = d
}

func GetDialect(name string) (d Dialect, ok bool) {
	d, ok = dialectMaps[strings.ToLower(name)]
	return
}