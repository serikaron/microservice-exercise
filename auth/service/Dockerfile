FROM serikaron/grpc-go as development

COPY . /auth
WORKDIR /auth

RUN go generate auth/service && \
    go build -o /auth-service auth/service

ENTRYPOINT /auth-service --auth-service-port $auth_service_port --auth-internal-service-port $auth_internal_service_port

