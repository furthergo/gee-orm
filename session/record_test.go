package session

import (
	"reflect"
	"testing"
)

func TestRecordInsert(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text, Age integer);").Exec()

	u1 := User{
		Name: "Tom",
		Age: 11,
	}

	u2 := User{
		Name: "Jack",
	}
	s.Insert(u1, u2)

	rows, _ := s.Raw("SELECT * FROM User").QueryRows()
	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.Name, &u.Age)
		if err == nil {
			users = append(users, u)
		}
	}
	if !reflect.DeepEqual([]User{u1, u2}, users) {
		t.Fatal("Insert record failed")
	}
}

func TestRecordFind(t *testing.T) {
	s := NewSession()
	var users []User
	err := s.Find(&users)
	if err != nil {
		t.Fatal("find record failed", err)
	}
	if len(users) != 2 {
		t.Fatal("error record count")
	}
}
