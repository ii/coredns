package federation

import (
	"fmt"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/middleware"
	"github.com/miekg/dns"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("federation", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	fed, err := federationParse(c)
	if err != nil {
		return middleware.Error("federation", err)
	}
	dnsserver.GetConfig(c).AddMiddleware(func(next middleware.Handler) middleware.Handler {
		fed.Next = next
		return nil
	})

	return nil
}

func federationParse(c *caddy.Controller) (Federation, error) {
	fed := New()

	for c.Next() {
		// federation [zones..]
		origins := make([]string, len(c.ServerBlockKeys))
		copy(origins, c.ServerBlockKeys)

		for c.NextBlock() {
			x := c.Val()
			switch c.Val() {
			case "fallthrough":
				fed.Fallthrough = true

			default:
				args := c.RemainingArgs()
				if len(args) != 1 {
					return fed, fmt.Errorf("need two arguments for federation: %q", args)
				}
				fed.f[x] = dns.Fqdn(args[0])
			}
		}

		for i := range origins {
			origins[i] = middleware.Host(origins[i]).Normalize()
		}

		fed.zones = origins

		if len(fed.f) == 0 {
			return fed, fmt.Errorf("at least one name to zone federation expected")
		}

		return fed, nil
	}

	return fed, nil
}
