FROM serikaron/grpc-go as development

COPY . /integration-test-client
WORKDIR /integration-test-client

RUN go generate auth/integration-test-client && \
    go build -o /auth-test-client auth/integration-test-client

FROM alpine:3.11 as release
COPY --from=development /auth-test-client /
ENTRYPOINT /auth-test-client --auth-service-port $auth_service_port --auth-internal-service-port $auth_internal_service_port
