package main

import (
	"context"
	"flag"

	"github.com/golang/glog"

	"passenger-tracker/pkg/server"
)

var (
	addr    = flag.String("addr", ":9090", "endpoint of the gRPC service")
	network = flag.String("network", "tcp", "a valid network type which is consistent to -addr")
)

func main() {
	flag.Parse()
	defer glog.Flush()

	ctx := context.Background()
	if err := server.Run(ctx, *network, *addr, ":8089"); err != nil {
		glog.Fatal(err)
	}
}
