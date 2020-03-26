FROM serika-registry:443/go-dev as development

COPY pkg /mse/pkg
COPY proto/*.proto proto/gen.go /mse/proto/
COPY auth /mse/auth

RUN go generate mse/proto && \
    go build -o /auth-service mse/auth/service

FROM alpine:3.11 as release
COPY --from=development /auth-service /
ENTRYPOINT /auth-service --auth-service-port $auth_service_port --auth-internal-service-port $auth_internal_service_port

