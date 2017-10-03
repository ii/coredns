# forward

*forward* facilitates proxying DNS messages to upstream resolvers.

The *forward* plugin is generally faster then *proxy* as it re-uses already openened sockets to
upstream. It supports UDP and TCP and uses inband healthchecking that is enabled by default.

## Syntax

In its most basic form, a simple forwarder uses this syntax:

~~~
forward FROM TO...
~~~

* **FROM** is the base domain to match for the request to be forwared.
* **TO...** are the destination endpoints to forward to.

By default health checks are performed every 0.5s. After two failed checks the upstream is
considered unhealthy. The health checks use a a non recursive DNS query (`. IN NS`) to get upstream
health. Any reponse that indicates the server is active (REFUSED, NOTIMPL, NXDOMAIN, SUCCESS) is
taken as a healthy upstream.

It uses a fixed buffer size of 4096 bytes for UDP packets.

Extra knobs are available with an expanded syntax:

~~~
proxy FROM TO... {
    except IGNORED_NAMES...
    force_tcp
    health_check [DURATION]
}
~~~

* **FROM** and **TO...** as above.
* **IGNORED_NAMES** in `except` is a space-separated list of domains to exclude from proxying.
  Requests that match none of these names will be passed through.
* `force_tcp`, use TCP even when the request comes in over UDP.
* `health_checks`, use a different **DURATION** for health checking, the default duration is 500ms.

## Metrics

TODO

## Examples

Proxy all requests within example.org. to a backend system:

~~~ corefile
. {
    forward example.org 127.0.0.1:9005
}
~~~

Load-balance all requests between three backends:

~~~ corefile
. {
    forward . 10.0.0.10:53 10.0.0.11:1053 10.0.0.12
}
~~~

Forward everything except requests to `miek.nl` or `example.org`

~~~ corefile
. {
    forward . 10.0.0.10:1234 {
        except miek.nl example.org
    }
}
~~~

Proxy everything except `example.org` using the host's `resolv.conf`'s nameservers:

~~~ corefile
. {
    proxy . /etc/resolv.conf {
        except miek.nl example.org
    }
}
~~~
