package nodecache

import (
	"strings"
	"time"

	"github.com/miekg/dns"
	"go.uber.org/zap"
)

func (p *NodeCache) QueryUpstream(domain string, typ uint16) *Answer {
	msg := new(dns.Msg)
	msg.SetQuestion(domain, typ)

	answer := new(Answer)

	for _, upstream := range p.upAddrs {
		if i := strings.Index(upstream, ":"); i < 0 {
			upstream += ":53"
		}

		dnsClient := dns.Client{}
		resp, rrt, err := dnsClient.Exchange(msg, upstream)
		if err != nil {
			zlog.Error("query_upstream_fail",
				zap.String("domain", domain),
				zap.String("upstream", upstream),
			)
			continue
		}

		if len(resp.Answer) == 0 {
			continue
		}

		answer.Msg = *resp
		answer.Expire = time.Now().Add(rrt)

		return answer
	}

	return nil
}
