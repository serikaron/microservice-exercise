FROM serikaron/grpc-go:latest as development

RUN mkdir /chat
WORKDIR /chat

COPY go.mod /chat
RUN go mod download

COPY proto /chat/proto
COPY pkg /chat/pkg
COPY monitor /chat/monitor

RUN go generate chat/monitor && \
    go build -o /monitor chat/monitor

FROM alpine:3.11 as release
COPY --from=development /monitor /
ENTRYPOINT /monitor --chat-service-port $chat_service_port

