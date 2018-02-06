package client

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/coredns/coredns/pb"
	"github.com/coredns/coredns/plugin/pkg/tls"

	"github.com/miekg/dns"
	"google.golang.org/grpc"
	creds "google.golang.org/grpc/credentials"
)

type Client struct {
	pbClient pb.DnsServiceClient
}

type Watch struct {
	WatchId int64
	Msgs chan *dns.Msg
	stream pb.DnsService_WatchClient
	client *Client
}

func NewClient(endpoint, cert, key, ca string) (*Client, error) {
	var tlsargs []string
	if cert != "" {
		tlsargs = append(tlsargs, cert)
	}

	if key != "" {
		tlsargs = append(tlsargs, key)
	}

	if ca != "" {
		tlsargs = append(tlsargs, ca)
	}

	var dialOpts []grpc.DialOption
	if len(tlsargs) == 0 {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	} else {
		tlsConfig, err := tls.NewTLSConfigFromArgs(tlsargs...)
		if err != nil {
			return nil, err
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds.NewTLS(tlsConfig)))
	}
	conn, err := grpc.Dial(endpoint, dialOpts...)
	if err != nil {
		return nil, err
	}
	return &Client{pbClient: pb.NewDnsServiceClient(conn)}, nil
}

func (c *Client) Query(req *dns.Msg) (*dns.Msg, error) {
        msg, err := req.Pack()
        if err != nil {
                return nil, err
        }

        reply, err := c.pbClient.Query(context.Background(), &pb.DnsPacket{Msg: msg})
        if err != nil {
                return nil, err
        }
        d := new(dns.Msg)
        err = d.Unpack(reply.Msg)
        if err != nil {
                return nil, err
        }
        return d, nil
}

func (c *Client) QueryNameAndType(qname string, qtype uint16) (*dns.Msg, error) {
	m := &dns.Msg{}
	m.SetQuestion(dns.Fqdn(qname), qtype)

	return c.Query(m)
}

func (c *Client) Watch(req *dns.Msg) (*Watch, error) {
	p, err := req.Pack()
	if err != nil {
		return nil, err
	}

	query := &pb.DnsPacket{Msg: p}
	cr := &pb.WatchRequest_CreateRequest{CreateRequest: &pb.WatchCreateRequest{Query: query}}

	stream, err := c.pbClient.Watch(context.Background())
	if err != nil {
		return nil, err
	}

	if err = stream.Send(&pb.WatchRequest{RequestUnion: cr}); err != nil {
		return nil, err
	}

	in, err := stream.Recv()
	if err == io.EOF {
		return nil, fmt.Errorf("Server returned EOF after attempt to create watch.")
	}
	if !in.Created {
		return nil, fmt.Errorf("Unexpected non-created response from server: %v", in)
	}
	w := &Watch{WatchId: in.WatchId, Msgs: make(chan *dns.Msg), stream: stream, client: c}
	go func() {
		for {
			in, err := w.stream.Recv()
			if err == io.EOF {
				close(w.Msgs)
				return
			}
			if err != nil {
				log.Printf("[ERROR] Watch %d failed to receive from gRPC stream: %s\n", w.WatchId, err)
				close(w.Msgs)
				return
			}

			if in.Created {
				log.Printf("[ERROR] Watch %d unexpected created response from server: %v\n", w.WatchId, in)
				close(w.Msgs)
				return
			}
			if in.Canceled {
				close(w.Msgs)
				return
			}

			r, err := w.client.Query(req)
			if err != nil {
				log.Printf("[ERROR] Error querying for changes: %s\n", err)
				close(w.Msgs)
				return
			}
			w.Msgs <- r
		}
	}()

	return w, nil
}

func (c *Client) WatchNameAndType(qname string, qtype uint16) (*Watch, error) {
	m := &dns.Msg{}
	m.SetQuestion(dns.Fqdn(qname), qtype)

	return c.Watch(m)
}


func (w *Watch) Stop() error {
	cr := &pb.WatchRequest_CancelRequest{CancelRequest: &pb.WatchCancelRequest{WatchId: w.WatchId}}
	return w.stream.Send(&pb.WatchRequest{RequestUnion: cr})
}
