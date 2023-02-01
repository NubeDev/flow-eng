package cmap

import "reflect"

type MapUtil struct {
	m map[string]interface{}
}

type MapUtils []MapUtil

func New(m map[string]interface{}) MapUtil {
	return MapUtil{m: m}
}

func (t MapUtil) child(key string) MapUtil {
	if t.m == nil {
		return MapUtil{}
	}

	v, ok := t.GetByValue(key, map[string]interface{}{})
	if !ok {
		return MapUtil{}
	}

	return MapUtil{m: v.(map[string]interface{})}
}

func (t MapUtil) GetAll() map[string]interface{} {
	return t.m
}

// func (t MapUtil) GetAll(keys ...string) MapUtil {
//	for _, k := range keys {
//		t = t.child(k)
//	}
//	return t
// }

// Key represents either the index of an array (int) or the key in a map (string).
type Key struct {
	v interface{}
}

func (k Key) IsArray() bool {
	_, ok := k.v.(int)
	return ok
}

func (k Key) IsMap() bool {
	_, ok := k.v.(string)
	return ok
}

func (k Key) Array() int {
	i, ok := k.v.(int)
	if !ok {
		panic("Array() called on non-array key")
	}
	return i
}

func (k Key) Map() string {
	i, ok := k.v.(string)
	if !ok {
		panic("Map() called on non-string key")
	}
	return i
}

type SetFunc func(k Key, value interface{}) (interface{}, bool)

// SetAll traverses all []interface{} and map[string]interface{} types and calls the
// fn (SetFunc) for each key/value pair. If the SetFunc for a given key/value pair
// returns boolean(true) as its 2nd return value, then said value will be updated to
// whatever SetFunc returned as the interface{}.
func (t MapUtil) SetAll(fn SetFunc) int {
	if t.m == nil {
		return 0
	}

	return setAll(t.m, fn)
}

func setAll(i interface{}, fn SetFunc) int {
	o := 0
	switch t := i.(type) {
	case map[string]interface{}:
		for k, v := range t {
			res, changed := fn(Key{k}, v)
			if changed {
				t[k] = res
				o++
			}

			o += setAll(v, fn)
		}
	case []interface{}:
		for k := range t {
			switch t[k].(type) {
			case map[string]interface{}:
				o += setAll(t[k], fn)
			default:
				res, changed := fn(Key{k}, t[k])
				if changed {
					t[k] = res
					o++
				}
			}
		}
	}
	return o
}

func (t MapUtil) FindAllWithKey(key string) MapUtils {
	if t.m == nil {
		return nil
	}

	return findAllWithKey(t.m, key)
}

func findAllWithKey(i interface{}, key string) MapUtils {
	o := MapUtils{}
	switch t := i.(type) {
	case map[string]interface{}:
		for k, v := range t {
			if k == key {
				o = append(o, MapUtil{t})
			}

			o = append(o, findAllWithKey(v, key)...)
		}
	case []interface{}:
		for _, v := range t {
			o = append(o, findAllWithKey(v, key)...)
		}
	}
	return o
}

func (t MapUtil) GetByValue(key string, i interface{}) (interface{}, bool) {
	if t.m == nil {
		return nil, false
	}

	v, ok := t.m[key]
	if !ok {
		return nil, false
	}

	if reflect.TypeOf(i) == reflect.TypeOf(v) {
		return v, true
	}

	return nil, false
}

func (t MapUtil) GetByKey(key string) (interface{}, bool) {
	if t.m == nil {
		return nil, false
	}
	v, ok := t.m[key]
	if !ok {
		return nil, false
	}
	return v, false
}

func (t MapUtil) Add(key string, i interface{}) bool {
	if t.m == nil {
		return false
	}
	t.m[key] = i
	return true
}

func (t MapUtil) Delete(key string) bool {
	if t.m == nil {
		return false
	}
	delete(t.m, key)
	return true
}
