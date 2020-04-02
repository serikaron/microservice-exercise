FROM serika-registry:443/go-dev as build

COPY pkg /mse/pkg
COPY proto/*.proto proto/gen.go /mse/proto/
COPY test/integration_test/chat /mse/test

RUN go generate mse/proto && \
    go test -c -o /chat-test mse/test

FROM alpine:3.11
COPY res /mse/res
COPY --from=build /chat-test /
ENTRYPOINT /chat-test --chat-port=$chat_port --cert-path $cert_path