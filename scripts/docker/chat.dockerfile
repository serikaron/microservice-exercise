FROM serika-registry:443/go-dev as development

COPY pkg /mse/pkg
COPY proto/*.proto proto/gen.go /mse/proto/
COPY chat /mse/chat

RUN go generate mse/proto && \
    go build -o /chat-service mse/chat/service

FROM alpine:3.11 as release
COPY --from=development /chat-service /
ENTRYPOINT /chat-service --chat-service-port $chat_service_port --chat-internal-service-port $chat_internal_service_port

