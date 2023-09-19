package mymem

import (
	"bytes"
	"fmt"
	"strings"
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

func (r *Rows) readAllItems() error {
	for r.Next() {
		if r.err != nil {
			return r.err
		}
	}
	return r.err
}

func (r *Rows) readValue() error {
	if r.conn == nil {
		return fmt.Errorf("closed network connection")
	}
	var err error
	var line []byte
	buf := make([]byte, 1)
	for {
		_, err = r.conn.Read(buf)
		if err != nil {
			return err
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
			return err
		}
		// read body
		buf = make([]byte, r.lenBody+2)
		_, err = r.conn.Read(buf)
		if err != nil {
			return err
		}
		r.values = strings.Split(strings.TrimSpace(string(buf)), r.delimiter)
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
	return ErrUnexpected
}
