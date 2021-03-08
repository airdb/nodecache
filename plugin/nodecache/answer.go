package nodecache

import (
	"github.com/miekg/dns"
	"time"
)

type Answer struct {
	Msg    dns.Msg
	Expire time.Time
	Entries []string
}

type Entry struct {
	Hdr HDR
	A string `json:"-"`
}

type HDR struct {
	Name string
	Rrtype uint16
	Class uint16
	Ttl uint32
	Rdlength uint16
}