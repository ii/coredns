FROM golang:alpine
MAINTAINER "Denver Williams <denver@ii.coop>"
ENV KUBECTL_VERSION=v1.5.2
ENV ARC=amd64

# Install Deps
RUN apk update
RUN apk add --update git bash util-linux wget tar curl build-base
RUN go get -u github.com/coredns/coredns/test

#Install Kubectl
RUN wget -O /usr/local/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/$ARC/kubectl && \
chmod +x /usr/local/bin/kubectl

#Copy Files

COPY files/ /coredns/e2e/

WORKDIR /go/src/github.com/coredns/coredns/test

CMD ["/bin/bash"]
