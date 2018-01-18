package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/coredns/coredns/pb"
	"github.com/coredns/coredns/plugin/pkg/tls"

	"github.com/miekg/dns"
	"google.golang.org/grpc"
	creds "google.golang.org/grpc/credentials"
)

func main() {
	var (
		endpoint string
		qname    string
		qtype    string
		verbose  bool
		insecure bool
		cert     string
		key      string
		ca       string
	)

	flag.BoolVar(&verbose, "v", false, "Log verbosely")
	flag.BoolVar(&insecure, "insecure", false, "Don't use TLS")
	flag.StringVar(&cert, "cert", "", "TLS cert PEM file path")
	flag.StringVar(&key, "key", "", "TLS key PEM file path")
	flag.StringVar(&ca, "ca", "", "TLS ca cert PEM file path")
	flag.StringVar(&endpoint, "server", "localhost:5553", "server host:port")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 || len(args) > 2 {
		panic(fmt.Errorf("Usage: dns-watch [flags] <name> [type]"))
	}
	qname = args[0]
	qtype = "A"
	if len(args) > 1 {
		qtype = args[1]
	}

	if verbose {
		fmt.Printf("Watching %s records for %s on server %s\n", qtype, qname, endpoint)
	}

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
	if insecure {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	} else {
		tlsConfig, err := tls.NewTLSConfigFromArgs(tlsargs...)
		if err != nil {
			panic(err)
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds.NewTLS(tlsConfig)))
	}
	conn, err := grpc.Dial(endpoint, dialOpts...)
	client := pb.NewWatchServiceClient(conn)

	stream, err := client.Watch(context.Background())
	if err != nil {
		panic(err)
	}

	m := &dns.Msg{}
	m.SetQuestion(dns.Fqdn(qname), dns.StringToType[qtype])

	p, err := m.Pack()
	if err != nil {
		panic(err)
	}
	query := &pb.DnsPacket{Msg: p}
	cr := &pb.WatchRequest_CreateRequest{CreateRequest: &pb.WatchCreateRequest{Query: query}}
	if err = stream.Send(&pb.WatchRequest{RequestUnion: cr}); err != nil {
		panic(err)
	}

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
		log.Printf("Got %v", in)
	}
}
