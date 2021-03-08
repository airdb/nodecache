package main

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"

	_ "airdb.io/airdb/nodecache/plugin/nodecache"
	_ "github.com/coredns/coredns/core/plugin"
)

func init() {
	dnsserver.Directives = append([]string{"nodecache"}, dnsserver.Directives...)
}

func main() {
	coremain.Run()
}
