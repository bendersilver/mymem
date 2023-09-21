package mymem

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"
)

func (r *Rows) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

func (r *Rows) Err() error {
	return r.err
}

func (r *Rows) Next() bool {
	r.doStep = false
	r.err = r.readValue()
	if r.err != nil {
		return false
	}
	return r.doStep
}

func (r *Rows) newConn(m *MySQLMemcached) *Rows {
	r.delimiter = m.delimiter
	r.conn, r.err = net.DialTimeout("tcp", m.host, time.Second)
	if r.err != nil {
		r.err = fmt.Errorf("mymem: connection error: %v", r.err)
	}
	return r
}

func (r *Rows) getContainers(name string) *Rows {
	if r.err != nil {
		return r
	}
	_, r.err = r.conn.Write([]byte(fmt.Sprintf("get @@containers.%s\r\n", name)))
	if r.err != nil {
		r.err = fmt.Errorf("mymem: write query get containers: %v", r.err)
		return r
	}
	r.readAllItems()
	if r.err != nil {
		if r.err == ErrNotFound {
			r.err = fmt.Errorf("mymem: required adding innodb_memcache.containers to memcache")
		}
		return r
	}
	if len(r.values) < 2 {
		r.err = fmt.Errorf("mymem: empty value_columns in table 'containers'")
		return r
	}
	r.container = &struct {
		Value []string
		Key   string
	}{r.values, r.values[0]}
	r.values = nil
	return r
}

func (r *Rows) cdNamespace(name string) *Rows {
	if r.err != nil {
		return r
	}
	_, r.err = r.conn.Write([]byte(fmt.Sprintf("get @@%s\r\n", name)))
	if r.err != nil {
		r.err = fmt.Errorf("mymem: write query cd namespace: %v", r.err)
		return r
	}
	r.readAllItems()
	r.values = nil
	return r
	// failed to locate entry in config table 'containers' in database 'innodb_memcache'
}

func (r *Rows) writeCMD(cmd string) *Rows {
	if r.err != nil {
		return r
	}
	_, r.err = r.conn.Write([]byte(cmd))
	if r.err != nil {
		r.err = fmt.Errorf("mymem: write cmd %s: %v", cmd, r.err)
	}
	return r
	// failed to locate entry in config table 'containers' in database 'innodb_memcache'
}

func (r *Rows) writeEndLine() (err error) {
	if r.err != nil {
		return r.err
	}
	var out []byte
	buf := make([]byte, 1)
	for {
		_, err = r.conn.Read(buf)
		if err != nil {
			return fmt.Errorf("mymem: read end line error: %v", err)
		}
		out = append(out, buf...)
		if bytes.HasSuffix(out, []byte("\r\n")) {
			break
		}
	}
	out = bytes.TrimSpace(out)
	switch string(out) {
	case "NOT_STORED":
		return ErrNotStored
	case "NOT_FOUND":
		return ErrNotFound
	case "ERROR":
		return ErrQuery
	case "DELETED", "STORED":
		return nil
	}
	return fmt.Errorf("mymem: unexpected line in response: '%s'", out)
}

func (r *Rows) readAllItems() {
	for r.Next() {
		if r.err != nil {
			return
		}
	}
}

func (r *Rows) readValue() error {
	if r.err != nil {
		return r.err
	}

	var err error
	var line []byte
	buf := make([]byte, 1)
	for {
		_, err = r.conn.Read(buf)
		if err != nil {
			return fmt.Errorf("mymem: read first line error: %v", err)
		}
		line = append(line, buf...)
		if bytes.HasSuffix(line, []byte("\r\n")) {
			break
		}
	}
	line = bytes.TrimSpace(line)
	if len(line) >= 5 && string(line[:5]) == "VALUE" {
		r.values = nil
		_, err = fmt.Sscanf(string(line), "VALUE %s %d %d\r\n", &r.Key, &r.Flag, &r.lenBody)
		if err != nil {
			return fmt.Errorf("mymem: ssca line error: %v", err)
		}

		// read body
		r.lenBody += 2
		line = nil

		var n int
		for r.lenBody > 0 {
			buf = make([]byte, r.lenBody)
			n, err = r.conn.Read(buf)
			if err != nil {
				return fmt.Errorf("mymem: read body error: %v", err)
			}
			r.lenBody -= n
			line = append(line, buf[:n]...)
		}
		r.original = string(line)
		r.values = strings.Split(strings.TrimSpace(string(line)), r.delimiter)
		r.doStep = true
		return nil

	} else if string(line) == "END" {
		if r.values == nil {
			return ErrNotFound

		}
		return nil

	} else if string(line) == "ERROR" {
		return ErrQuery

	}
	return fmt.Errorf("mymem: unexpected line in response: '%s'", line)
}
