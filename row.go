package mymem

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
)

func (r *Rows) Values() ([]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	return append([]string{r.Key}, r.values...), nil
}

func (r *Rows) Map() (map[string]string, error) {
	vals, err := r.Values()
	if err != nil {
		return nil, err
	} else if r.container == nil {
		return nil, ErrContainers
	}
	m := make(map[string]string)
	for i, k := range r.container.Value {
		m[k] = vals[i]
	}
	return m, nil
}

func (r *Rows) Scan(ptrs ...any) error {
	vals, err := r.Values()
	if err != nil {
		return err
	} else if len(ptrs) != len(vals) {
		return ErrCount
	}
	for i, ptr := range ptrs {
		if ptr == nil {
			continue
		}
		err = r.setValue(vals[i], ptr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Rows) setValue(val string, ptr any) (err error) {
	tp := reflect.TypeOf(ptr)
	if tp.Kind() != reflect.Pointer {
		return ErrNonPrt
	}
	tp = tp.Elem()
	var b []byte
	switch ptr.(type) {
	case B64String, B64Bool, B64Int, B64Uint, B64Float:
		b, err = base64.StdEncoding.DecodeString(val)
	}
	switch tp.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Map:
		b, err = base64.StdEncoding.DecodeString(val)
	}
	if b != nil {
		switch tp.Kind() {
		case reflect.Struct, reflect.Map:
			return json.Unmarshal(b, ptr)
		case reflect.Slice:
			if tp.Elem().Kind() != reflect.Uint8 {
				reflect.ValueOf(ptr).Elem().SetBytes(b)
			} else {
				err = json.Unmarshal(b, ptr)
			}
			return
		case reflect.String:
			reflect.ValueOf(ptr).Elem().SetString(string(b))
			return
		default:
			val = string(b)
		}

	}
	_, err = fmt.Sscan(val, ptr)
	if err != nil {
		return fmt.Errorf("failed to assign value: %v", err)
	}
	return
}

func (r *Rows) ScanStruct(ptr any) error {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Pointer {
		return ErrNonPrt
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return ErrNonStruct
	}
	vals, err := r.Values()
	if err != nil {
		return err
	}
	t := reflect.TypeOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		jname, ok := t.Field(i).Tag.Lookup("json")
		if !ok {
			continue
		}
		for ix, key := range r.container.Value {
			if key == jname {
				err = r.setValue(vals[ix], v.Field(i).Addr().Interface())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
