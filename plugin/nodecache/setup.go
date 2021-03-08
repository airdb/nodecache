package nodecache

import (
	"strconv"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"go.uber.org/zap"
)

func init() {
	plugin.Register("nodecache", setup)
}

var zlog *zap.Logger

func setup(c *caddy.Controller) error {
	zlog = InitLogger()

	nodeCache, err := nodeCacheParse(c)
	if err != nil {
		return plugin.Error("nodecache", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		nodeCache.Next = next
		return nodeCache
	})

	zlog.Info("nodecache_start")
	return nil
}

func nodeCacheParse(c *caddy.Controller) (*NodeCache, error) {
	params := make(map[string]string)

	for c.Next() {
		args := c.RemainingArgs()
		for _, arg := range args {
			splits := strings.SplitN(arg, "=", 2)
			params[splits[0]] = splits[1]
		}
	}

	cachePath := params["cache_path"]
	upAddrs := strings.Split(params["up_addrs"], ",")
	cacheSize, _ := strconv.Atoi(params["cache_size"])

	return NewNodeCache(cachePath, cacheSize, upAddrs), nil
}
