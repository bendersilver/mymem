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
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Pointer {
		return ErrNonPrt
	}
	typeName := v.Elem().Type().Name()
	switch typeName {
	case "B64Byte", "B64String", "B64Bool", "B64Int", "B64Uint", "B64Float":
		b, e := base64.StdEncoding.DecodeString(val)
		if e != nil {
			return fmt.Errorf("failed to decode base64 value: %v", e)
		}
		val = string(b)

	}
	_, ok := reflect.TypeOf(ptr).MethodByName("B64UnmarshalJSON")
	if ok {
		b, e := base64.StdEncoding.DecodeString(val)
		if e != nil {
			return fmt.Errorf("failed to decode base64 value: %v", e)
		}
		return json.Unmarshal(b, ptr)
	}

	switch v.Elem().Interface().(type) {
	case string, B64String:
		v.Elem().SetString(val)
		return
	case []byte, B64Byte:
		v.Elem().SetBytes([]byte(val))
	default:
		_, err = fmt.Sscan(val, ptr)
		if err != nil {
			return fmt.Errorf("failed to assign value: %v", err)
		}
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
