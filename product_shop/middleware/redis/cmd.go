package redis

func Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	return rdb.Eval(script, keys, args).Result()
}
