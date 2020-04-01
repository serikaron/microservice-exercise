FROM serika-registry:443/go-dev as development

COPY pkg /mse/pkg
COPY proto/*.proto proto/gen.go /mse/proto/
COPY auth /mse/auth

RUN go generate mse/proto && \
    go build -o /auth mse/auth

FROM alpine:3.11 as release
COPY --from=development /auth /
ENTRYPOINT /auth --auth-host "" --auth-port $auth_port

