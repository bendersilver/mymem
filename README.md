# requery
```sql
INSERT INTO innodb_memcache.containers VALUES('containers', 'innodb_memcache', 'containers', 'name', 'key_columns|value_columns', 0, 0, 0, 'PRIMARY');
```

# custom type
B64* - convert from base64 to type in Scan and ScanStruct method

QueryRow
```go
   
func TestQueryRow(t *testing.T) {
	// Table: dash_session

	// Columns:
	// 	k		varchar(256) PK
	// 	vb64	text
	// 	ttl		bigint


    // NewMySQLMemcached(<host:port>, <delimiter>)
	m := NewMySQLMemcached("127.0.0.1:11211", "|")

	type b64 struct {
		B64Struct `json:"-"`
		String    string
		Int       int
		Bool      bool
	}
	b, err := json.MarshalIndent(&b64{String: "string value", Int: 100500}, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s", b)                   // {
                                        //         "String": "string value",
                                        //         "Int": 100500,
                                        //         "Bool": false
                                        // }

	b64val := base64.StdEncoding.EncodeToString(b)
	t.Logf("\n%s", b64val)              // ewoJIlN0cmluZyI6ICJzdHJpbmcgdmFsdWUiLAoJIkludCI6IDEwMDUwMCwKCSJCb29sIjogZmFsc2UKfQ==

	err = m.Set("dash_session ", "unique_key", b64val, 6400)
	if err != nil {
		t.Fatal(err)
	}
	row := m.QueryRow("dash_session", "unique_key")

	mp, err := row.Map()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n", mp)                     // map[k:unique_key ttl:6400 vb64:ewoJIlN0cmluZyI6ICJzdHJpbmcgdmFsdWUiLAoJIkludCI6IDEwMDUwMCwKCSJCb29sIjogZmFsc2UKfQ==]

	vals, err := row.Values()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n", vals)                   // [unique_key ewoJIlN0cmluZyI6ICJzdHJpbmcgdmFsdWUiLAoJIkludCI6IDEwMDUwMCwKCSJCb29sIjogZmFsc2UKfQ== 6400]


	var example struct {
		K   string `json:"k"`
		Val b64    `json:"vb64"`
		TTL int64  `json:"ttl"`
	}


	row.ScanStruct(&example)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%v", example)             // {unique_key {{} string value 100500 false} 6400}

	var k string
	var val b64
	var ttl int64
	err = row.Scan(&k, &val, &ttl)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n", k, val, ttl)            // unique_key {{} string value 100500 false} 6400

	var valStr B64String
	err = row.Scan(nil, &valStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n", valStr)                 // {
                                        //         "String": "string value",
                                        //         "Int": 100500,
                                        //         "Bool": false
                                        // }

	var valByte B64Byte
	err = row.Scan(nil, &valByte, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s", valByte)             // {
                                        //         "String": "string value",
                                        //         "Int": 100500,
                                        //         "Bool": false
                                        // }
}

```


Query
```go
func TestQuery(t *testing.T) {
	m := NewMySQLMemcached("127.0.0.1:11211", "|")
	// https://dev.mysql.com/doc/refman/8.0/en/innodb-memcached-multiple-get-range-query.html
	// To get all values greater than B, enter get @>B:
	// get @>B

	// To get all values less than M, enter get @<M:
	// get @<M

	// To get all values less than and including M, enter get @<=M:
	// get @<=M

	// To get values greater than B but less than M, enter get @>B@<M:
	// get @>B@<M

	rows, err := m.Query("dash_session ", "@>")
	if err != nil {
		t.Fatal(err)
	}
    defer rows.Close()
	var k string
	for rows.Next() {
		err = rows.Scan(&k, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("\n", k)
	}
    // output
    //     query_test.go:: 
    //         1lLi496744z1JggxWTwNgyxs9HF9C1YRUAXWiK2hi6bQd9k3PV4riC7ZlDVT7xnq
    //     query_test.go:: 
    //         4lnXsUujHrZ0HVKp1k4O1t4IrmmwCoC0mIHlmA4HKlbX3fzwSLcQNYJFfPOt33oT
    //     query_test.go:: 
    //         4VjcfLL8ISorgGGRczcWhITIIkSXSoTKym3N7MFwgLz1nvPKhP3yyDCWdNKUrXsi
    //     query_test.go:: 
    //         CtyHPJ4QZFakVrdCiSE2IRfKQN8dxrA9sr7fKT8Q7uPvFc0yPDoeZnqDRflZIfNL
    //     query_test.go:: 
    //         GkhCvBxDpaRArLYTaL0xb8d5ndDMQWAfCGKRwSjzYjTDDgq5MgooYeMDgQIor1Gm
    //     query_test.go:: 
    //         ytTaz5upXazVa1aWhp7JAxcbjdUeWXJstt8zHMyG8nWf45gclEBzbMgx14RKgrwr
	if rows.Err() != nil {
		t.Fatal(err)
	}
}

```