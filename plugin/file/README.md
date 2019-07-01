# file

## Name

*file* - enables serving zone data from an RFC 1035-style master file.

## Description

The *file* plugin is used for an "old-style" DNS server. It serves from a preloaded file that exists
on disk.

If the zone file contains signatures (i.e., is signed using DNSSEC), correct DNSSEC answers are
returned. Only NSEC is supported! If you use this setup *you* are responsible for re-signing the
zonefile.

However the plugin can also automatically resign the zone for you in this setup you need to create
a CSK key (multiple ones are supported) and an output directory where CoreDNS can write the resigned
zone. Thus CoreDNS will:

* Resign the zone with the CSK (the ZSK/KZK split is *not* supported) starting every Thursday at
  15:00 UTC.
* Create signatures that have an inception of -3H and expiration of +3W for every key given.
* Add replace *all* CDS with the keys given.
* Update the SOA's serial number to the Unix epoch of when the signing happens. This will overwrite
  the previous serial number.

Keys are named (following BIND9): `K<name>+<alg>+<id>.key` and `K<name>+<alg>+<id>.private`. The keys
must not be included in your zone; they will be added by CoreDNS when the zone is signed. These keys
can be generated with `coredns-keygen` or BIND9's `dnssec-keygen`.

## Syntax

~~~
file DBFILE [ZONES...]
~~~

* **DBFILE** the database file to read and parse. If the path is relative, the path from the *root*
  directive will be prepended to it.
* **ZONES** zones it should be authoritative for. If empty, the zones from the configuration block
    are used.

If you want to round-robin A and AAAA responses look at the *loadbalance* plugin.

~~~
file DBFILE [ZONES... ] {
    transfer to ADDRESS...
    reload DURATION
    upstream
    dnssec KEYDIR [DIR]
}
~~~

* `transfer` enables zone transfers. It may be specified multiples times. `To` or `from` signals
  the direction. **ADDRESS** must be denoted in CIDR notation (e.g., 127.0.0.1/32) or just as plain
  addresses. The special wildcard `*` means: the entire internet (only valid for 'transfer to').
  When an address is specified a notify message will be sent whenever the zone is reloaded.
* `reload` interval to perform a reload of the zone if the SOA version changes. Default is one minute.
  Value of `0` means to not scan for changes and reload. For example, `30s` checks the zonefile every 30 seconds
  and reloads the zone when serial changes.
* `upstream` resolve external names found (think CNAMEs) pointing to external names. This is only
  really useful when CoreDNS is configured as a proxy; for normal authoritative serving you don't
  need *or* want to use this. CoreDNS will resolve CNAMEs against itself.
* `directory` specifies the directory where CoreDNS should save zones that are being signed. If not
  given this defaults to `/var/lib/coredns`. This setting is only used if `dnssec` is given.
* `dnssec` enables DNSSEC zone signing for all zones specified. **KEYDIR** is used to read the keys
  from. The signed zones are saved to **DIR** from the `directory` option, which defaults to
  `/var/lib/coredns` when not given. The zones are saved under the name `Z<name>.signed.` (the `Z`
  is added to not hide the root zone).

## Examples

Load the `example.org` zone from `example.org.signed` and allow transfers to the internet, but send
notifies to 10.240.1.1

~~~ corefile
example.org {
    file example.org.signed {
        transfer to *
        transfer to 10.240.1.1
    }
}
~~~

Or use a single zone file for multiple zones:

~~~
. {
    file example.org.signed example.org example.net {
        transfer to *
        transfer to 10.240.1.1
    }
}
~~~

Note that if you have a configuration like the following you may run into a problem of the origin
not being correctly recognized:

~~~
. {
    file db.example.org
}
~~~

We omit the origin for the file `db.example.org`, so this references the zone in the server block,
which, in this case, is the root zone. Any contents of `db.example.org` will then read with that
origin set; this may or may not do what you want.
It's better to be explicit here and specify the correct origin. This can be done in two ways:

~~~
. {
    file db.example.org example.org
}
~~~

Or

~~~
example.org {
    file db.example.org
}
~~~

## Also See

The DNSSEC RFC: RFC 4033, RFC 4034 and RFC 4035, coredns-keygen(1) and dnssec-keygen(8).
