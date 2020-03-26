FROM serikaron/grpc-go:latest as development

RUN mkdir /chat
WORKDIR /chat

COPY go.mod /chat
RUN go mod download

COPY proto /chat/proto
COPY pkg /chat/pkg
COPY chater /chat/chater

RUN go generate chat/chater && \
    go build -o /chater chat/chater

FROM alpine:3.11 as release
COPY --from=development /chater /
ENTRYPOINT /chater --chat-service-port $chat_service_port

