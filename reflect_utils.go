package redis

import "reflect"

const (
	RedisFiledTag = "redis"
)

func ReflectSaveEntityArgs(rev reflect.Value) []interface{} {
	var args []interface{}
	ret := rev.Type()

	for i := 0; i < rev.NumField(); i++ {
		revF := rev.Field(i)
		if revF.Kind() == reflect.Struct {
			args = ReflectSaveEntityArgs(revF)
			continue
		}

		retF := ret.Field(i)
		fn, ok := retF.Tag.Lookup(RedisFiledTag)
		if ok {
			args = append(args, fn, revF.Interface())
		}
	}

	return args
}
