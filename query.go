package mymem

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Query -
func (m *MySQLMemcached) Query(name string, key any) (*Rows, error) {
	name = strings.TrimSpace(name)
	var r Rows
	defer func() {
		if r.err != nil {
			r.Close()
		}
	}()
	r.newConn(m).
		getContainers(name).
		cdNamespace(name).
		writeCMD(fmt.Sprintf("get %v\r\n", key))
	return &r, nil
}

// QueryRow -
func (m *MySQLMemcached) QueryRow(name string, key any) *Rows {
	name = strings.TrimSpace(name)
	var r Rows
	defer r.Close()
	r.newConn(m).
		getContainers(name).
		writeCMD(fmt.Sprintf("get @@%s.%v\r\n", name, key))
	if r.err != nil {
		return &r
	}
	r.Next()
	r.Key = fmt.Sprintf("%v", key)
	return &r
}

func (m *MySQLMemcached) Exists(name string, key any) (ok bool, err error) {
	_, err = m.QueryRow(name, key).Values()
	if err == ErrNotFound {
		return false, nil
	} else if err == nil {
		return true, nil
	}
	return
}

// Delete -
func (m *MySQLMemcached) Delete(name string, key any) error {
	var r Rows
	defer r.Close()
	r.newConn(m).
		cdNamespace(strings.TrimSpace(name)).
		writeCMD(fmt.Sprintf("delete %v\r\n", key))
	if r.err != nil {
		return r.err
	}
	return r.writeEndLine()
}

// Set -
func (m *MySQLMemcached) Set(name string, key any, args ...any) error {
	var r Rows
	defer r.Close()
	var values []string
	for _, v := range args {
		if v == nil {
			return fmt.Errorf("nil value not supported")
		}
		tp := reflect.TypeOf(v)
		if tp.Kind() != reflect.Pointer {
			v = reflect.ValueOf(v).Elem().Interface()
			tp = tp.Elem()
		}
		switch v.(type) {
		case B64String, B64Bool, B64Int, B64Uint, B64Float:
			values = append(values, base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", v))))
			continue
		}

		var b []byte
		var ok bool

		switch tp.Kind() {
		case reflect.Struct, reflect.Map, reflect.Slice:
			if b, ok = v.([]byte); !ok {
				b, r.err = json.Marshal(v)
			}
			values = append(values, base64.StdEncoding.EncodeToString(b))
			continue
		}

		values = append(values, fmt.Sprintf("%v", v))
	}
	body := strings.Join(values, m.delimiter)

	r.newConn(m).
		writeCMD(fmt.Sprintf("set @@%s.%v 0 0 %d\r\n%s\r\n", strings.TrimSpace(name), key, len(body), body))
	if r.err != nil {
		return r.err
	}
	return r.writeEndLine()
}
