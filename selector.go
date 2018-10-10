package gomemcache

import "net"

type ServerSelector interface {
	PickServer(key string)(net.Addr, error)
	Each(func(net.Addr) error) error
}