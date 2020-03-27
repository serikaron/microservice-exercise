FROM serika-registry:443/go-dev as build

COPY pkg /mse/pkg
COPY proto /mse/proto
COPY testing /mse/testing

RUN go generate mse/proto && \
    go test -c -o /chat-test mse/testing/chat

FROM alpine:3.11
COPY --from=build /chat-test /
ENTRYPOINT /chat-test --chat-service-port=$chat_service_port