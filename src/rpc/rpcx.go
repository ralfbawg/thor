package rpc

import (
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"context"
	"flag"
)

var (
	addr = flag.String("addr", "localhost:8972", "server address")
)

func PeerToPeer(host string, port string, serviceName string, args []string) {
	d := client.NewPeer2PeerDiscovery("tcp@" + *addr, "")
}

func RcpxServer() {

	s := server.NewServer()
	//s.Register(new(Arith), "")
	s.RegisterName("Arith", new(Arith), "")
	err := s.Serve("tcp", *addr)
	if err != nil {
		panic(err)
	}
}

type Arith struct{}

// the second parameter is not a pointer
func (t *Arith) Mul(ctx context.Context, args example.Args, reply *example.Reply) error {
	reply.C = args.A * args.B
	return nil
}

func getServiceAddr(host string, port string, serviceName string) string {
	return host + ":" + port
}
