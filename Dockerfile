FROM golang:1.9.2

RUN go get -u github.com/golang/dep/cmd/dep

RUN mkdir -p /go/src/github.com/sample-crd

COPY . /go/src/github.com/sample-crd
WORKDIR /go/src/github.com/sample-crd

RUN dep ensure -vendor-only -v
