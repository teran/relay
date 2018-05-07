FROM golang

ADD . /go/src/github.com/teran/mail-relay
RUN cd /go/src/github.com/teran/mail-relay && CGO_ENABLED=0 go build -o bin/relay .

FROM scratch

COPY --from=0 /go/src/github.com/teran/mail-relay/bin/relay /relay

ENTRYPOINT ["/relay"]
