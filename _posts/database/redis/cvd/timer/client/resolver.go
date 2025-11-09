package client

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

const scheme = "grpc"

type builder struct{}

func (*builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &simpleResolver{
		target: target,
		cc:     cc,
	}
	addrs := strings.Split(target.Endpoint(), ",")
	r.addrs = make([]resolver.Address, 0, len(addrs))
	for _, addr := range addrs {
		r.addrs = append(r.addrs, resolver.Address{
			Addr: addr,
		})
	}
	r.start()
	return r, nil
}

func (*builder) Scheme() string {
	return scheme
}

type simpleResolver struct {
	target resolver.Target
	addrs  []resolver.Address
	cc     resolver.ClientConn
}

func (r *simpleResolver) start() {
	r.cc.UpdateState(resolver.State{Addresses: r.addrs})
}

func (*simpleResolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (*simpleResolver) Close() {}

func init() {
	resolver.Register(&builder{})
}
