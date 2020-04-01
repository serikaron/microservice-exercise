FROM serika-registry:443/go-dev as development

COPY pkg /mse/pkg
COPY proto/*.proto proto/gen.go /mse/proto/
COPY test/integration_test/auth /mse/test

RUN go generate mse/proto && \
    go test -c -o /auth-test mse/test

FROM alpine:3.11 as release
COPY --from=development /auth-test /
ENTRYPOINT /auth-test --auth-port $auth_port