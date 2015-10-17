package beard

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Data struct {
	Value interface{}

	// valueOf is a cache of value's reflect.Value
	valueOf *reflect.Value
}

func (d *Data) Get(k string) *Data {
	if k == "." {
		// dot notations just returns itself
		return d
	}
	if v := getValue(k, d.Value); v != nil {
		return &Data{
			Value: v,
		}
	}

	return nil
}

func (d *Data) Len() int {
	if d.IsSlice() {
		return d.ValueOf().Len()
	}

	return 0
}

func (d *Data) IsSlice() bool {
	return d.ValueOf().Kind() == reflect.Slice
}

func (d *Data) ValueOf() *reflect.Value {
	if d.valueOf == nil {
		v := reflect.ValueOf(d.Value)

		d.valueOf = &v
	}

	return d.valueOf
}

func (d *Data) Index(n int) *Data {
	return &Data{
		Value: d.ValueOf().Index(n).Interface(),
	}
}

func (d *Data) Bytes() []byte {
	var b []byte
	switch t := d.Value.(type) {
	case string:
		return []byte(t)
	case int:
		return strconv.AppendInt(b, int64(t), 10)
	case int64:
		return strconv.AppendInt(b, t, 10)
	case float64:
		return strconv.AppendFloat(b, float64(t), 'G', -1, 64)
	case bool:
		return strconv.AppendBool(b, t)
	case []byte:
		return t
	case reflect.Value:
		return []byte(t.String())

	default:
		return []byte(fmt.Sprintf("%s", t))
	}

	return nil
}

// getValue finds the value of k within source.
// This source can represent a tree like structure can must be traversed to find
// the value of the full path of k. `k` can be a json path, eg. `foo.bar`
func getValue(k string, source interface{}) interface{} {
	tr, br := splitpath(k)
	// TODO handle a is blank

	var v reflect.Value

	v, ok := source.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(source)
	}

	switch v.Kind() {
	case reflect.Map:
		v = v.MapIndex(reflect.ValueOf(tr))
	case reflect.Struct:
		v = v.FieldByName(tr)
	case reflect.Ptr:
		return getValue(k, v.Elem())
	default:
		return nil
	}

	if !v.IsValid() {
		return nil
	}
	var i interface{} = v

	if v.CanInterface() {
		i = v.Interface()
	}
	if br != "" {
		return getValue(br, i)
	}

	return i
}

const pathDelim = '.'

func splitpath(path string) (string, string) {
	i := strings.IndexByte(path, pathDelim)
	if i != -1 {
		return path[:i], path[i+1:]
	}

	return path, ""
}
