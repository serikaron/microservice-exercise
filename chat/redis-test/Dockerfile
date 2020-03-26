FROM serikaron/grpc-go

RUN mkdir /chat
WORKDIR /chat

COPY go.mod /chat
RUN go mod download

COPY proto /chat/proto
COPY pkg /chat/pkg
COPY redis-test /chat/redis-test

RUN go generate chat/redis-test && \
    go build -o /redis-test chat/redis-test

ENTRYPOINT /redis-test --redis-port $redis_port
