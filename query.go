package mymem

import (
	"fmt"
	"net"
	"time"
)

func (m *MySQLMemcached) Query(name string, key any) (*Rows, error) {
	var r Rows
	r.delimiter = m.delimiter
	var err error
	defer func() {
		if err != nil {
			r.Close()
		}
	}()

	r.conn, err = net.DialTimeout("tcp", m.host, time.Second)
	if err != nil {
		return nil, fmt.Errorf("mmcd: connection error: %v", err)
	}

	// get containers table
	_, err = r.conn.Write([]byte(fmt.Sprintf("get @@containers.%s\r\n", name)))
	if err != nil {
		return nil, fmt.Errorf("mmcd: write query: %v", err)
	}
	err = r.readAllItems()
	if err != nil {
		if err == ErrNotFound {
			return nil, fmt.Errorf("required adding innodb_memcache.containers to memcache")
		}
	}
	if len(r.values) < 2 {
		return nil, fmt.Errorf("mmcd: empty value_columns in table 'containers'")
	}
	r.container = &struct {
		Value []string
		Key   string
	}{r.values, r.values[0]}

	_, err = r.conn.Write([]byte(fmt.Sprintf("get @@%s\r\n", name)))
	if err != nil {
		return nil, fmt.Errorf("mmcd: write query: %v", err)
	}
	err = r.readAllItems()
	// failed to locate entry in config table 'containers' in database 'innodb_memcache'
	if err != nil {
		return nil, err
	}

	_, err = r.conn.Write([]byte(fmt.Sprintf("get %v\r\n", key)))
	if err != nil {
		return nil, fmt.Errorf("mmcd: write query: %v", err)
	}

	return &r, nil
}

func (m *MySQLMemcached) QueryRow(name string, key any) (r *Rows) {
	var err error
	r = new(Rows)
	r.delimiter = m.delimiter
	defer r.Close()

	r.conn, err = net.DialTimeout("tcp", m.host, time.Second)
	if err != nil {
		r.err = fmt.Errorf("mmcd: connection error: %v", err)
		return
	}

	// get containers table
	_, err = r.conn.Write([]byte(fmt.Sprintf("get @@containers.%s\r\n", name)))
	if err != nil {
		r.err = fmt.Errorf("mmcd: write query: %v", err)
		return
	}
	err = r.readAllItems()
	if err == nil {
		if len(r.values) < 2 {
			r.err = fmt.Errorf("mmcd: empty value_columns in table 'containers'")
			return
		}
		r.container = &struct {
			Value []string
			Key   string
		}{r.values, r.values[0]}
	}

	_, err = r.conn.Write([]byte(fmt.Sprintf("get @@%s.%v\r\n", name, key)))
	if err != nil {
		r.err = fmt.Errorf("mmcd: write query: %v", err)
		return
	}
	r.Next()
	return
}
