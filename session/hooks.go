package session

import (
	"github.com/futhergo/gee-orm/log"
	"reflect"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

func (s *Session)CallMethod(name string, value interface{}) {
	method := reflect.ValueOf(s.RefTable().Model).MethodByName(name)
	if value != nil {
		method = reflect.ValueOf(value).MethodByName(name)
	}

	if method.IsValid() {
		params := []reflect.Value{reflect.ValueOf(s)}
		if v := method.Call(params); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
	return
}