package dnstap

// Mirror some of the tap. constants so we can reference them from this plugin
// and not the vendored directory meaning out-of-tree middleware can use this.

import (
	tap "github.com/dnstap/golang-dnstap"
)

const (
	// Mirror constants from the golang-dnstap package.
	SocketProtocolTCP        = tap.SocketProtocol_TCP
	SocketProtocolUDP        = tap.SocketProtocol_UDP
	MessageForwarderQuery    = tap.Message_FORWARDER_QUERY
	MessageForwarderResponse = tap.Message_FORWARDER_RESPONSE
)
