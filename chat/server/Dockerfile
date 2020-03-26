FROM serikaron/grpc-go:latest as development

RUN mkdir /chat
WORKDIR /chat

COPY go.mod /chat
RUN go mod download

COPY proto /chat/proto
COPY pkg /chat/pkg
COPY server /chat/server

RUN go generate chat/server && \
    go build -o /chat-service chat/server

FROM alpine:3.11 as release
COPY --from=development /chat-service /
ENTRYPOINT /chat-service --chat-service-port $chat_service_port

