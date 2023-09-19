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
		if r.conn != nil && err != nil {
			r.conn.Close()
		}
	}()

	r.conn, err = net.DialTimeout("tcp", m.host, time.Second)
	if err != nil {
		return nil, fmt.Errorf("mmcd: connection error: %v", err)
	}
	if m.containers {
		_, err = r.conn.Write([]byte(fmt.Sprintf("get @@containers.%s\r\n", name)))
		err = r.readAllItems()
		if err != nil {
			return nil, ErrFailed
		}
	}

	_, err = r.conn.Write([]byte(fmt.Sprintf("get @@%s\r\n", name)))
	if err != nil {
		return nil, fmt.Errorf("mmcd: write query: %v", err)
	}
	err = r.readAllItems()
	// failed to locate entry in config table 'containers' in database 'innodb_memcache'
	if err != nil {
		return nil, err
	}

	return &r, nil
}
