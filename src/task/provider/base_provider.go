package provider
type providerI interface {
	GetData() interface{}
}
type BaseProvider struct {
	providerI
}
