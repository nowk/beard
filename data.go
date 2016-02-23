package beard

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Data struct {
	Value interface{}

	block *block

	k          string
	as         string
	isKeyValue bool
	i          int
	keys       []reflect.Value

	// valueOf is a cache of value's reflect.Value
	valueOf *reflect.Value
}

func (d *Data) As(as ...string) {
	d.isKeyValue = false // ensure false toggle

	switch len(as) {
	case 0:
		return
	case 1:
		d.as = as[0]
	case 2:
		d.k = as[0]
		d.as = as[1]

		d.isKeyValue = true
	default:
		// TODO handle it should not be allowed to be set more than 2 as vars
	}
}

func (d *Data) Get(k string) *Data {
	if d.isKeyValue && d.k != "" && d.k == k {
		return &Data{
			Value: d.getKey(d.k),
		}
	}

	if d.isKeyValue && d.as != "" && d.as == k {
		switch v := d.getKey(d.k).(type) {
		case reflect.Value:
			k = v.String()
		case string:
			k = v
		default:
			// TODO handle
		}
	}

	// dot notations just returns itself
	if k == "." || (d.as != "" && d.as == k) {
		return d
	}

	if v := d.getValue(k, d.Value); v != nil {
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
	if n, ok := d.IsKeys(); ok {
		return n
	}

	return 1
}

func (d *Data) IsSlice() bool {
	if d.Value == nil {
		return false
	}

	return d.ValueOf().Kind() == reflect.Slice
}

func (d *Data) IsKeys() (int, bool) {
	if !d.isKeyValue {
		return 0, false
	}
	if d.Value == nil {
		return 0, false
	}
	v, ok := d.Value.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(d.Value)
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		return len(v.MapKeys()), true
	case reflect.Struct:
		return v.NumField(), true
	}

	return 0, false
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

var keyName = func(r1, r2 *reflect.Value) bool {
	return r1.String() < r2.String()
}

type by func(*reflect.Value, *reflect.Value) bool

func (b by) Sort(r []reflect.Value) {
	s := &keySorter{
		keys: r,
		by:   b,
	}

	sort.Sort(s)
}

type keySorter struct {
	keys []reflect.Value
	by   by
}

func (k *keySorter) Len() int {
	return len(k.keys)
}

func (k *keySorter) Swap(i, j int) {
	k.keys[i], k.keys[j] = k.keys[j], k.keys[i]
}

func (k *keySorter) Less(i, j int) bool {
	return k.by(&k.keys[i], &k.keys[j])
}

func (d *Data) getKey(k string) interface{} {
	// look at the parent's data set to figure out the type
	val, ok := d.block.data.Value.(reflect.Value)
	if !ok {
		val = reflect.ValueOf(d.block.data.Value)
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var v interface{}

	switch val.Kind() {
	case reflect.Map:
		// save keys to struct and access keys via saved value, else we run into
		// issues with map keys not maintaining order.
		if d.keys == nil {
			keys := val.MapKeys()

			by(keyName).Sort(keys)

			d.keys = keys
		}

		v = d.keys[d.i]
	case reflect.Struct:
		f := val.Type().Field(d.i)

		v = f.Name

	default:
		// we only find keys on maps and structs, all else will return the index
		v = d.i
	}

	return v
}

// getValue finds the value of the path within source.
// The path can be represented as a json path, eg a.b.c and will traverse the
// source to find said path.
func (d *Data) getValue(path string, source interface{}) interface{} {
	i := strings.Index(path, d.as+".")
	if i == 0 {
		path = path[len(d.as)+1:]
	}
	tr, br := splitpath(path)
	if tr == "" {
		return nil
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
		return d.getValue(br, inf)
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
