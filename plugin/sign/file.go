package sign

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/coredns/coredns/plugin/file"
	"github.com/miekg/dns"
)

const outputDirectory = "/var/lib/coredns"

// write writes out the zone file to a temporary file which can then be moved into place.
func write() error {
	tmpFile, err := ioutil.TempFile(outputDirectory, "signed-")
	if err != nil {
		return err
	}
	defer tmpFile.Close()
	// write it out

	return nil
}

// parse parses the zone in filename and returns a new Zone or an error. This
// is similar to the Parse function in the *file* plugin. However On parse the
// record type RRSIG, DNSKEY, CDNSKEY and CDS are *not* included in the
// returned zone.
func parse(f io.Reader, origin, fileName string) (*file.Zone, error) {
	zp := dns.NewZoneParser(f, dns.Fqdn(origin), fileName)
	zp.SetIncludeAllowed(true)
	z := file.NewZone(origin, fileName)
	seenSOA := false

	for rr, ok := zp.Next(); ok; rr, ok = zp.Next() {
		if err := zp.Err(); err != nil {
			return nil, err
		}

		switch rr.(type) {
		case *dns.SOA:
			seenSOA = true
		case *dns.RRSIG, *dns.DNSKEY, *dns.CDNSKEY, *dns.CDS:
			// drop
		default:
			if err := z.Insert(rr); err != nil {
				return nil, err
			}
		}
	}
	if !seenSOA {
		return nil, fmt.Errorf("file %q has no SOA record", fileName)
	}

	return z, nil
}
