package beard

import (
	"reflect"
	"strings"
)

// getvof finds the value of k within source.
// This source can represent a tree like structure can must be traversed to find
// the value of the full path of k. `k` can be a json path, eg. `foo.bar`
func getvof(k string, source interface{}) interface{} {
	tr, br := splitpath(k)
	// TODO handle a is blank
	// TODO handle a is a dot (.)

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
		return getvof(k, v.Elem())
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
		return getvof(br, i)
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
