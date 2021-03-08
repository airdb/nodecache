package nodecache

import (
	"context"
	"go.uber.org/zap"
	"strings"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// Validation is a plugin that validates incoming DNS queries
type NodeCache struct {
	Next plugin.Handler

	upAddrs   []string
	cachePath string
	cacheSize int
	fastcache *fastcache.Cache
}

const defaultCachePath = "dns-cache"

// ServeDNS implements the plugin.Handler interface.
func (p NodeCache) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	domain := state.Name()
	dnsType := state.QType()

	if p.fastcache.Has([]byte(domain)) {
		// Check cache.
		rr, err := p.getCache(domain)
		if err != nil {
		}

		msg := new(dns.Msg)
		for _, r := range rr {
			msg.Answer = append(msg.Answer, *r)
		}
		if msg != nil {
			msg.SetReply(r)

			err = w.WriteMsg(msg)
			if err != nil {
			}
		}
	} else {
		answer := p.QueryUpstream(domain, dnsType)

		answer.Msg.SetReply(r)

		err := w.WriteMsg(&answer.Msg)
		if err != nil {
			return dns.RcodeNameError, err
		}

		p.saveCache(domain, answer)
		return dns.RcodeSuccess, nil
	}

	return dns.RcodeSuccess, nil
}

// Name implements the Handler interface.
func (p NodeCache) Name() string { return "nodecache" }

func (p *NodeCache) saveCache(domain string, m *Answer) {
	var entries []string

	for i := range m.Msg.Answer {
		m.Msg.Answer[i].Header().Ttl = 30
		entries = append(entries, m.Msg.Answer[i].String())
	}

	value := strings.Join(entries, ";")

	key := []byte(domain)
	p.fastcache.Set(key, []byte(value))
}

func (p *NodeCache) getCache(domain string) ([]*dns.RR, error) {
	key := []byte(domain)
	v := p.fastcache.Get(nil, key)

	entries := strings.Split(string(v), ";")

	var aa []*dns.RR
	for _, entry := range entries {
		rr, err := dns.NewRR(entry)
		if err != nil {
			zlog.Warn("new_rr_fail", zap.String("entry", entry))
			continue
		}

		aa = append(aa,  &rr)
	}

	return aa, nil
}

func NewNodeCache(cachePath string, cacheSize int, upAddrs []string) *NodeCache {
	p := &NodeCache{
		upAddrs:   upAddrs,
		cachePath: cachePath,
		cacheSize: cacheSize,
	}

	if p.cachePath == "" {
		p.cachePath = defaultCachePath
	}

	// Set default cache size 5000.
	if p.cacheSize <= 0 {
		p.cacheSize = 5000
	}

	p.fastcache = fastcache.New(p.cacheSize)

	go p.FlushCache()

	return p
}

func (p *NodeCache) FlushCache() {
	for {
		err := p.fastcache.SaveToFile(p.cachePath)
		if err != nil {
		}

		time.Sleep(time.Minute)
	}
}
