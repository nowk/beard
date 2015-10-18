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
	// dot notations just returns itself
	if k == "." {
		return d
	}
	if v := getValue(k, d.Value); v != nil {
		return &Data{
			Value: v,
		}
	}

	return nil
}

// Len returns the length of the data object. Any non nil object that is not a
// slice will be returned with a value of 1
func (d *Data) Len() int {
	if d.Value == nil {
		return 0
	}
	if d.IsSlice() {
		return d.ValueOf().Len()
	}

	return 1
}

func (d *Data) IsSlice() bool {
	if d.Value == nil {
		return false
	}

	return d.ValueOf().Kind() == reflect.Slice
}

func (d *Data) ValueOf() *reflect.Value {
	if d.Value == nil {
		return nil
	}
	if d.valueOf == nil {
		v := reflect.ValueOf(d.Value)

		d.valueOf = &v
	}

	return d.valueOf
}

func (d *Data) Index(n int) *Data {
	if d.Len() == 0 {
		return nil
	}

	v := d.ValueOf().Index(n).Interface()

	return &Data{Value: v}
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

// getValue finds the value of the path within source.
// The path can be represented as a json path, eg a.b.c and will traverse the
// source to find said path.
func getValue(path string, source interface{}) interface{} {
	tr, br := splitpath(path)
	if tr == "" {
		// TODO handle a is blank
	}

	v, ok := source.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(source)
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		v = v.MapIndex(reflect.ValueOf(tr))
	case reflect.Struct:
		v = v.FieldByName(tr)

	default:
		return nil
	}

	if !v.IsValid() {
		return nil
	}

	var inf interface{} = v

	if v.CanInterface() {
		inf = v.Interface()
	}

	if br != "" {
		return getValue(br, inf)
	}

	return inf
}

const pathDelim = '.'

func splitpath(path string) (string, string) {
	i := strings.IndexByte(path, pathDelim)
	if i != -1 {
		return path[:i], path[i+1:]
	}

	return path, ""
}
