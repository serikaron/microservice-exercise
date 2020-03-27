FROM serika-registry:443/go-dev as development

COPY proto /mse/proto
COPY pkg /mse/pkg
COPY chat/monitor /mse/chat/monitor

RUN go generate mse/proto && \
    go build -o /monitor mse/chat/monitor

FROM alpine:3.11 as release
COPY --from=development /monitor /
ENTRYPOINT /monitor --chat-service-port $chat_service_port