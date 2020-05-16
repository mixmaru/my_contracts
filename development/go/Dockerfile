FROM golang:1.14.2

RUN go get github.com/go-delve/delve/cmd/dlv
RUN go get -v github.com/rubenv/sql-migrate/...

WORKDIR /root/my_contracts
