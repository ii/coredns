# Configuration

There are to parts to configuration CoreDNS. The first is determining which plugins you want to
compile into CoreDNS. The binaries we provide have all plugins as listed in `plugin.cfg` compiled
in. Adding or removing is [easy](/link/to/howto), but shouldn't normally be done by end users.

Thus most users use the *Corefile* to configure CoreDNS. When CoreDNS starts, and the `-conf` flag is
not given it will look for a file named `Corefile` in the current directory. That files consists out
of one or more Server stanzas. Each Server stanza lists one or more Plugins. Those Plugins may be
further configured with Directives. The ordering of the Plugin in the Corefile *does not determine*
the order of the plugin chain.

As said (/link) the plugin chain ordering is fixed and determined via plugin.cfg during the
compilation phase.

## Server Stanza

Server grouping


## Plugins

Directives


## External Plugin

Compile time enabled.

