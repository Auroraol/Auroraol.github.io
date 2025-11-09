package client

import (
	"gitlab.xiaoduoai.com/golib/grpc_resolver/dns"
	"google.golang.org/grpc/resolver"
)

func init() {
	resolver.Register(dns.NewBuilder(-1))
}
