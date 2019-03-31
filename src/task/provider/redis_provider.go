package provider

import (
	"datasource"
	"reflect"
	"errors"
)

type RedisProvider struct {
	redis datasource.RedisDatasource
	BaseProvider
}

func (redis *RedisProvider) GetData(key interface{}) (interface{}, error) {
	value := reflect.ValueOf(key)
	if value.Kind() != reflect.String {
		return nil, errors.New("gob: attempt to decode into a non-pointer")
	}
	return nil, nil
}
