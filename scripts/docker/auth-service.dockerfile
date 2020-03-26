FROM serikaron/go-dev as development

COPY auth /mse/auth

RUN go generate mse/auth/service && \
    go build -o /auth-service mse/auth/service

FROM alpine:3.11 as release
COPY --from=development /auth-service /
ENTRYPOINT /auth-service --auth-service-port $auth_service_port --auth-internal-service-port $auth_internal_service_port

