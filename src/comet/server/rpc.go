package server

import "github.com/smallnest/rpcx/server"

type RpcServer struct {
	ServerConfig struct {
		port int
	}
}

func init() {
	s := server.NewServer()
	s.RegisterName("service", new(), "")
	s.Serve("tcp", addr)
}
