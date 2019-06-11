package rpc

type ServiceConfig struct {
	name   string
	host   string
	port   string
	method string
	desc   string
}

//func Init() {
//	Register(&ServiceConfig{})
//	Service()
//}

func Service() {

}

func Register(conf *ServiceConfig) bool {
	return true
}
