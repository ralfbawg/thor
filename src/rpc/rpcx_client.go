package rpc

import (
	"flag"
	"github.com/smallnest/rpcx/client"
	"context"
	"log"
)

var (
	addr = flag.String("addr", "localhost:8972", "server address")
)

func DoRpcxCall(host string, port string, serviceName, args []interface{}) {

	flag.Parse()

	d := client.NewPeer2PeerDiscovery("tcp@" + *addr, "")
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := example.Args{
		A: 10,
		B: 20,
	}

	reply := &example.Reply{}
	err := xclient.Call(context.Background(), "Mul", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Printf("%d * %d = %d", args.A, args.B, reply.C)

}
