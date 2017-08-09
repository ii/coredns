package kubernetes

import (
	"reflect"
	"testing"

	"github.com/miekg/dns"
)

func expectString(t *testing.T, function, qtype, query string, r *recordRequest, field, expected string) {
	ref := reflect.ValueOf(r)
	refField := reflect.Indirect(ref).FieldByName(field)
	got := refField.String()
	if got != expected {
		t.Errorf("Expected %v(%v, \"%v\") to get %v == \"%v\". Instead got \"%v\".", function, query, qtype, field, expected, got)
	}
}

func TestParseRequest(t *testing.T) {

	var tcs map[string]string

	k := Kubernetes{Zones: []string{"inter.webs.test"}}
	f := "parseRequest"

	// Test a valid SRV request
	//
	query := "_http._tcp.webs.mynamespace.svc.inter.webs.test."
	r, e := k.parseRequest(query, dns.TypeSRV)
	if e != nil {
		t.Errorf("Expected no error from parseRequest(%v, \"SRV\"). Instead got '%v'.", query, e)
	}

	tcs = map[string]string{
		"port":      "http",
		"protocol":  "tcp",
		"endpoint":  "",
		"service":   "webs",
		"namespace": "mynamespace",
		"podOrSvc":  Svc,
		"zone":      "inter.webs.test",
	}
	for field, expected := range tcs {
		expectString(t, f, "SRV", query, &r, field, expected)
	}

	// Test wildcard acceptance
	//
	query = "*.any.*.any.svc.inter.webs.test."
	r, e = k.parseRequest(query, dns.TypeSRV)
	if e != nil {
		t.Errorf("Expected no error from parseRequest(\"%v\", \"SRV\"). Instead got '%v'.", query, e)
	}

	tcs = map[string]string{
		"port":      "*",
		"protocol":  "any",
		"endpoint":  "",
		"service":   "*",
		"namespace": "any",
		"podOrSvc":  Svc,
		"zone":      "inter.webs.test",
	}
	for field, expected := range tcs {
		expectString(t, f, "SRV", query, &r, field, expected)
	}

	// Test A request of endpoint
	query = "1-2-3-4.webs.mynamespace.svc.inter.webs.test."
	r, e = k.parseRequest(query, dns.TypeA)
	if e != nil {
		t.Errorf("Expected no error from parseRequest(\"%v\", \"A\"). Instead got '%v'.", query, e)
	}
	tcs = map[string]string{
		"port":      "",
		"protocol":  "",
		"endpoint":  "1-2-3-4",
		"service":   "webs",
		"namespace": "mynamespace",
		"podOrSvc":  Svc,
		"zone":      "inter.webs.test",
	}
	for field, expected := range tcs {
		expectString(t, f, "A", query, &r, field, expected)
	}

	// Test NS request
	query = "inter.webs.test."
	r, e = k.parseRequest(query, dns.TypeNS)
	if e != nil {
		t.Errorf("Expected no error from parseRequest(\"%v\", \"NS\"). Instead got '%v'.", query, e)
	}
	tcs = map[string]string{
		"port":      "",
		"protocol":  "",
		"endpoint":  "",
		"service":   "",
		"namespace": "",
		"podOrSvc":  "",
		"zone":      "inter.webs.test",
	}
	for field, expected := range tcs {
		expectString(t, f, "NS", query, &r, field, expected)
	}

	// Invalid query tests
	invalidAQueries := []string{
		"_http._tcp.webs.mynamespace.svc.inter.webs.test.", // A requests cannot have port or protocol
		"servname.ns1.srv.inter.nets.test.",                // A requests must have zone that matches corefile

	}
	for _, q := range invalidAQueries {
		_, e = k.parseRequest(q, dns.TypeA)
		if e == nil {
			t.Errorf("Expected error from %v(\"%v\", \"A\").", f, q)
		}
	}

	invalidSRVQueries := []string{
		"_http._pcp.webs.mynamespace.svc.inter.webs.test.", // SRV protocol must be tcp or udp
		"_http._tcp.ep.webs.ns.svc.inter.webs.test.",       // SRV requests cannot have an endpoint
		"_*._*.webs.mynamespace.svc.inter.webs.test.",      // SRV request with invalid wildcards
		"_http._tcp",
		"_tcp.test.",
		".",
	}

	for _, q := range invalidSRVQueries {
		_, e = k.parseRequest(q, dns.TypeSRV)
		if e == nil {
			t.Errorf("Expected error from %v(\"%v\", \"SRV\").", f, q)
		}
	}
}
