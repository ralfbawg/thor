package provider

type providerI interface {
	GetData(key interface{}) (interface{}, error)
}
type BaseProvider struct {
	providerI
}
