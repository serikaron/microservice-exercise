FROM serikaron/go-dev as development

COPY chat /mse/chat

RUN go generate mse/chat/service && \
    go build -o /chat-service mse/chat/service

FROM alpine:3.11 as release
COPY --from=development /chat-service /
ENTRYPOINT /chat-service --chat-service-port $chat_service_port --chat-internal-service-port $chat_internal_service_port

