package go_template

import (
	"reflect"
	"sort"
	"strings"
)

func reflectObjListField(objList any, key string) (v reflect.Value, isMap bool, isMethod bool) {
	if objList == nil {
		return reflect.Value{}, false, false
	}

	v = reflect.ValueOf(objList)
	if v.Kind() != reflect.Slice {
		panic("objList is not slice")
	}
	if v.Len() == 0 {
		return v, false, false
	}
	elem := v.Index(0)
	if elem.Kind() == reflect.Interface {
		elem = elem.Elem()
	}
	if elem.Kind() == reflect.Map {
		isMap = true
	} else if m := elem.MethodByName(key); m.IsValid() {
		isMethod = true
	} else {
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		field := elem.FieldByName(key)
		if !field.IsValid() {
			panic("invalid key " + key)
		}
	}
	return v, isMap, isMethod
}

func reflectObjListFieldElemValue(v reflect.Value, idx int, isMap, isMethod bool, key string) reflect.Value {
	elem := v.Index(idx)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}

	if isMap {
		return elem.MapIndex(reflect.ValueOf(key))
	} else if isMethod {
		return elem.MethodByName(key).Call(nil)[0]
	} else {
		return elem.FieldByName(key)
	}
}

func sortBy(objList any, key string) any {
	v, isMap, isMethod := reflectObjListField(objList, key)
	if !v.IsValid() || v.Len() == 0 {
		return nil
	}
	sort.Slice(objList, func(i, j int) bool {
		valI := reflectObjListFieldElemValue(v, i, isMap, isMethod, key)
		valJ := reflectObjListFieldElemValue(v, j, isMap, isMethod, key)

		switch valI.Kind() {
		case reflect.String:
			return valI.String() < valJ.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return valI.Int() < valJ.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return valI.Uint() < valJ.Uint()
		case reflect.Float32, reflect.Float64:
			return valI.Float() < valJ.Float()
		default:
			return strings.Compare(valI.String(), valJ.String()) < 0
		}
	})
	return objList
}

func anyBy(objList any, key string, value any) bool {
	v, isMap, isMethod := reflectObjListField(objList, key)
	if !v.IsValid() || v.Len() == 0 {
		return false
	}
	for i := 0; i < v.Len(); i++ {
		val := reflectObjListFieldElemValue(v, i, isMap, isMethod, key)
		if val.IsValid() && val.Interface() == value {
			return true
		}
	}
	return false
}
