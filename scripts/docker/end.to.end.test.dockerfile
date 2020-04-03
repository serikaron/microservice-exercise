FROM serika-registry:443/go-dev as build

COPY pkg /mse/pkg
COPY proto/*.proto proto/gen.go /mse/proto/
COPY test/end_to_end /mse/test

RUN go generate mse/proto && \
    go test -c -o /end-to-end-test mse/test

FROM alpine:3.11
COPY res /mse/res
COPY --from=build /end-to-end-test /
ENTRYPOINT /end-to-end-test --chat-port=$chat_port --cert-path $cert_path --auth-port=$auth_port