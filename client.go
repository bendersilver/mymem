package mymem

import (
	"net"
)

type MySQLMemcached struct {
	host       string
	delimiter  string
	containers bool
}

func NewMySQLMemcached(addr, delimiter string, containers bool) *MySQLMemcached {
	return &MySQLMemcached{
		delimiter:  delimiter,
		host:       addr,
		containers: containers,
	}
}

type Rows struct {
	Key  string
	Flag int

	conn      net.Conn
	delimiter string
	container *struct {
		Value string
		Key   string
	}
	buf, line []byte
	lenBody   int

	values []string
	count  int
	doStep bool
	err    error
}
