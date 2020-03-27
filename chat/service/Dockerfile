FROM serika-registry:443/go-dev as development

COPY proto /mse/proto
COPY pkg /mse/pkg
COPY chat/service /mse/chat/service

RUN go generate mse/proto && \
    go build -o /chat-service mse/chat/service

FROM alpine:3.11 as release
COPY --from=development /chat-service /
ENTRYPOINT /chat-service --chat-service-port $chat_service_port

