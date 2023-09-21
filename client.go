package mymem

import (
	"net"
)

type B64String string
type B64Bool bool
type B64Int int64
type B64Uint uint64
type B64Float float64

// INSERT INTO innodb_memcache.containers
// VALUES('containers', 'innodb_memcache', 'containers', 'name', 'key_columns|value_columns', 0, 0, 0, 'PRIMARY');
type MySQLMemcached struct {
	host      string
	delimiter string
}

func NewMySQLMemcached(addr, delimiter string) *MySQLMemcached {
	return &MySQLMemcached{
		delimiter: delimiter,
		host:      addr,
	}
}

type Rows struct {
	Key  string
	Flag int

	conn      net.Conn
	delimiter string
	container *struct {
		Value []string
		Key   string
	}

	lenBody  int
	values   []string
	original string
	doStep   bool
	err      error
}
