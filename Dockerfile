FROM golang:1.12 AS builder

COPY . /app
WORKDIR /app

RUN go mod download
RUN go build -o bin/nova-dex-ctl -v -ldflags '-s -w' cli/admincli/main.go && \
  go build -o bin/adminapi -v -ldflags '-s -w' cli/adminapi/main.go && \
  go build -o bin/api -v -ldflags '-s -w' cli/api/main.go && \
  go build -o bin/engine -v -ldflags '-s -w' cli/engine/main.go && \
  go build -o bin/launcher -v -ldflags '-s -w' cli/launcher/main.go && \
  go build -o bin/watcher -v -ldflags '-s -w' cli/watcher/main.go && \
  go build -o bin/websocket -v -ldflags '-s -w' cli/websocket/main.go && \
  go build -o bin/maker -v -ldflags '-s -w' cli/maker/main.go

FROM alpine
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN apk update && \
  apk add sqlite ca-certificates && \
  rm -rf /var/cache/apk/*

COPY --from=builder /app/db /db/
COPY --from=builder /app/bin/* /bin/
