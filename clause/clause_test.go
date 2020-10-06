package clause

import (
	"reflect"
	"testing"
)

type user struct {
	Name string
	Age int
}

func RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	dt := destValue.Type()
	var fieldValues []interface{}
	for i := 0; i < dt.NumField(); i++ {
		fieldValues = append(fieldValues, destValue.FieldByIndex([]int{i}).Interface())
	}
	return fieldValues
}

func TestClause(t *testing.T) {
	c := Clause{}
	c.Set(INSERT, "User", []string{"Name", "Age"})
	//tom := user{Name: "Tom", Age: 13}
	//jack := user{Name: "Jack", Age: 29}
	//c.Set(VALUES, RecordValues(tom), RecordValues(jack))
	c.Set(VALUES, []interface{}{"Tom", 13}, []interface{}{"Jack"})
	q, vars := c.Build(INSERT, VALUES)
	if q != "INSERT INTO User (Name, Age) VALUES (?, ?), (?, ?)" {
		t.Fatal("build query string failed")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 13, "Jack"}) {
		t.Fatal("build args failed")
	}
}
