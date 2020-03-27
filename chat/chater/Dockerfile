FROM serika-registry:443/go-dev as development

COPY proto /mse/proto
COPY pkg /mse/pkg
COPY chat/chater /mse/chat/chater

RUN go generate mse/proto && \
    go build -o /chater mse/chat/chater

FROM alpine:3.11 as release
COPY res/certs/cert.pem /
COPY --from=development /chater /
ENTRYPOINT /chater --chat-service-port $chat_service_port

