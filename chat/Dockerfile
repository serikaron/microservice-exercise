FROM serikaron/grpc-go:latest as development

COPY . /chat/

WORKDIR /chat

RUN go generate chat/server

