package mymem

import (
	"testing"
)

func Test(t *testing.T) {
	var j struct {
		Name string    `json:"name"`
		Val  B64String `json:"bval"`
	}

	m := NewMySQLMemcached("127.0.0.1:11211", "|")
	rows, err := m.Query("conf", "mc.dash_login_users")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var name string
	var val B64String
	for rows.Next() {
		t.Log(rows.Map())
		t.Log(rows.Scan(&name, &val))
		t.Log(name, val)
		t.Log(rows.ScanStruct(&j))
		t.Log(j)
	}
	t.Fatal(rows.Err())
}

func TestQueryRow(t *testing.T) {

	var j struct {
		Name string `json:"name"`
		Val  struct {
			B64Struct
			Table string
			Key   string
			Value string
		} `json:"bval"`
	}
	// j.UnmarshalJSON([]byte(`{"name": "2"}`))
	// t.Fatal(json.Unmarshal([]byte(`{"name": "2"}`), &j), j)

	m := NewMySQLMemcached("127.0.0.1:11211", "|")
	err := m.QueryRow("conf", "mc.dash_login_users").ScanStruct(&j)
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(j)
}
