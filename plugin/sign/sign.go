package sign

import (
	"fmt"
	"os"
	"time"

	"github.com/coredns/coredns/plugin/file"
	"github.com/coredns/coredns/plugin/file/tree"
)

type Sign struct {
	keys       []Pair
	expiration time.Duration
	inception  time.Duration
	dbfile     string
}

func signFunc(e *tree.Elem) bool {
	for qtype, rrs := range e.m {
		println(qtype)
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
	// sign it
	z.Tree.Do(signFunc)

	// print it
	z.Tree.Do(func(e *tree.Elem) bool {
		for _, r := range e.All() {
			fmt.Println(r.String())
		}
		return false
	})

	return nil
}
