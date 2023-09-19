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
	r.values = nil
	var err error
	r.buf = make([]byte, 1)
	for {
		_, err = r.conn.Read(r.buf)
		if err != nil {
			return err
		}
		r.line = append(r.line, r.buf...)
		if bytes.HasSuffix(r.buf, []byte("\r\n")) {
			break
		}
	}
	r.line = bytes.TrimSpace(r.line)
	if len(r.line) >= 5 && string(r.line[:5]) == "VALUE" {
		_, err = fmt.Sscanf(string(r.line), "VALUE %s %d %d\r\n", &r.Key, &r.Flag, &r.lenBody)
		if err != nil {
			return err
		}
		// read body
		r.buf = make([]byte, r.lenBody)
		_, err = r.conn.Read(r.buf)
		if err != nil {
			return err
		}
		r.buf = bytes.TrimSpace(r.buf)
		r.values = strings.Split(strings.TrimSpace(string(r.buf)), r.delimiter)
		r.count++
		r.doStep = true
		return nil

	} else if string(r.line) == "END" {
		if r.count == 0 {
			return ErrNotFound

		}
		r.Close()
		return nil

	} else if string(r.line) == "ERROR" {
		return ErrQuery

	}
	return ErrUnexpected
}
