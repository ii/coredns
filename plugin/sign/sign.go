package sign

import (
	"fmt"
	"os"
	"time"

	"github.com/coredns/coredns/plugin/file"
	"github.com/coredns/coredns/plugin/file/tree"

	"github.com/miekg/dns"
)

type Sign struct {
	keys       []Pair
	expiration time.Duration
	inception  time.Duration
	dbfile     string
}

func (s Sign) signFunc(e *tree.Elem) bool {
	for qtype, rrs := range e.M() {
		println(qtype)
		println(rrs[0].String())
	}
	return false
}

func (s Sign) Sign(origin string) error {
	rd, err := os.Open(s.dbfile)
	if err != nil {
		return err
	}

	z, err := file.Parse(rd, origin, s.dbfile, 0)
	if err != nil {
		return err
	}
	for _, key := range s.keys {
		ds1 := key.Public.ToDS(dns.SHA1)
		ds2 := key.Public.ToDS(dns.SHA256)
		println(ds1.ToCDS().String())
		println(ds2.ToCDS().String())
		println(key.Public.ToCDNSKEY().String())
	}

	// sign it
	z.Tree.Do(s.signFunc)

	// print it
	z.Tree.Do(func(e *tree.Elem) bool {
		for _, r := range e.All() {
			fmt.Println(r.String())
		}
		return false
	})

	return nil
}
