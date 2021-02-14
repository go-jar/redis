package redis

import (
	"reflect"
)

const (
	FiledTag = "redis"
)

func ReflectSaveEntityArgs(rev reflect.Value) []interface{} {
	if rev.Kind() == reflect.Ptr {
		rev = rev.Elem()
	}

	var args []interface{}

	if rev.Kind() == reflect.Struct {
		ret := rev.Type()

		for i := 0; i < rev.NumField(); i++ {
			revF := rev.Field(i)
			if revF.Kind() == reflect.Struct || revF.Kind() == reflect.Ptr {
				tmpArgs := ReflectSaveEntityArgs(revF)
				args = append(args, tmpArgs...)
			} else {
				retF := ret.Field(i)
				fn, ok := retF.Tag.Lookup(FiledTag)
				if ok {
					args = append(args, fn, revF.Interface())
				}
			}
		}
	}

	return args
}
