FROM serika-registry:443/go-dev as development

COPY pkg /mse/pkg
COPY proto/*.proto proto/gen.go /mse/proto/
COPY chat /mse/chat

RUN go generate mse/proto && \
    go build -o /chat mse/chat

FROM alpine:3.11 as release
COPY res /mse/res
COPY --from=development /chat /
ENTRYPOINT /chat --chat-port $chat_port --chat-host "" --redis-port $redis_port --cert-path $cert_path

