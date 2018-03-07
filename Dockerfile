FROM alpine:latest

# only need ca-certificates & openssl if want to use DNS over TLS.
RUN apk --no-cache add bind-tools ca-certificates openssl && update-ca-certificates

ADD coredns /coredns

EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]
