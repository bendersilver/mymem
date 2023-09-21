package mymem

import (
	"testing"
)

func TestSet(t *testing.T) {
	// Table: dash_session

	// Columns:
	// 	k		varchar(256) PK
	// 	vb64	text
	// 	ttl		bigint

	m := NewMySQLMemcached(":11211", "|")

	type b64 struct {
		String string
		Int    int
		Bool   bool
	}

	err := m.Set("dash_session ", "unique_key", &b64{String: "string value", Int: 100500}, 6400)
	if err != nil {
		t.Fatal(err)
	}

	var example struct {
		K   string `json:"k"`
		Val b64    `json:"vb64"`
		TTL int64  `json:"ttl"`
	}

	row := m.QueryRow("dash_session", "unique_key")

	mp, err := row.Map()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n", mp)

	vals, err := row.Values()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n", vals)

	row.ScanStruct(&example)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%v", example)

	var k string
	var val b64
	var ttl int64
	err = row.Scan(&k, &val, &ttl)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n", k, val, ttl)

	var valStr B64String
	err = row.Scan(nil, &valStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n", valStr)

	var valByte []byte
	err = row.Scan(nil, &valByte, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s", valByte)
	err = m.Delete("dash_session", "unique_key")
	if err != nil {
		t.Fatal(err)
	}

	rows, err := m.Query("queue ", "@>")
	if err != nil {
		t.Fatal(err)
	}

	for rows.Next() {
		err = rows.Scan(&k, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("\n", k)
	}
	if rows.Err() != nil {
		t.Fatal(err)
	}
	rows.Close()

}

func TestQuery(t *testing.T) {

	m := NewMySQLMemcached(":11211", "|")
	// https://dev.mysql.com/doc/refman/8.0/en/innodb-memcached-multiple-get-range-query.html
	// To get all values greater than B, enter get @>B:
	// get @>B

	// To get all values less than M, enter get @<M:
	// get @<M

	// To get all values less than and including M, enter get @<=M:
	// get @<=M

	// To get values greater than B but less than M, enter get @>B@<M:
	// get @>B@<M

	rows, err := m.Query("partner_keys", 8)
	if err != nil {
		t.Fatal(err)
	}

	defer rows.Close()
	var k string
	for rows.Next() {
		err = rows.Scan(nil, &k)
		if err != nil {
			t.Fatal(err)
			continue
		}
		t.Log(k)

	}
	if rows.Err() != nil {
		t.Fatal(rows.Err())
	}
}
