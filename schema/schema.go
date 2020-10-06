package schema

import (
	"github.com/futhergo/gee-orm/dialect"
	"go/ast"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag string
}

type Schema struct {
	Model interface{}
	Name string
	Fields []*Field
	FieldsName []string
	FieldsMap map[string]*Field
}

func (s *Schema)GetField(name string) *Field {
	return s.FieldsMap[name]
}

func Parse(desc interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(desc)).Type()
	s := &Schema{
		Model: desc,
		Name: modelType.Name(),
		FieldsMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		sf := modelType.Field(i)
		if !sf.Anonymous && ast.IsExported(sf.Name) {
			f := &Field{
				Name: sf.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(sf.Type))),
			}
			if tag, ok := sf.Tag.Lookup("geeorm"); ok {
				f.Tag = tag
			}
			s.Fields = append(s.Fields, f)
			s.FieldsName = append(s.FieldsName, f.Name)
			s.FieldsMap[f.Name] = f
		}
	}
	return s
}

func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}