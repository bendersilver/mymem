package mymem

import "errors"

var (
	ErrNotFound   = errors.New("mymem: not found")
	ErrNotStored  = errors.New("mymem: not srored data")
	ErrTable      = errors.New("mymem: table not found")
	ErrMulti      = errors.New("mymem: result returned multiple records")
	ErrQuery      = errors.New("mymem: error query")
	ErrUnexpected = errors.New("mymem: unexpected line in get response")
	ErrFailed     = errors.New("mymem: failed to locate entry in config table 'containers' in database 'innodb_memcache'")
	ErrContainers = errors.New("mymem: config table 'containers' not in memcache")
	ErrCount      = errors.New("mymem: column count doesn't match value count at row")
	ErrNonPrt     = errors.New("mymem: non pointer argument")
	ErrNonStruct  = errors.New("mymem: non struct argument")
)
