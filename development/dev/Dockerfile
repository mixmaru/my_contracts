FROM golang:1.14.2
RUN apt update && apt -y upgrade

# デバッガ
RUN go get github.com/go-delve/delve/cmd/dlv
# DBマイグレーションツール
RUN go get -v github.com/rubenv/sql-migrate/...
# モックライブラリ
RUN go get github.com/golang/mock/mockgen

WORKDIR /root/my_contracts
