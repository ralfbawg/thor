package provider

import "datasource"

type RedisProvier struct {
	redis datasource.RedisDatasource
	BaseProvider
}

