FROM grpc-go:latest as development

WORKDIR $GOPATH/src/echo

COPY . $GOPATH/src/echo/

RUN mkdir -p $GOPATH/src/echo/proto && \
#    go generate && \
    protoc -Iproto --go_out=plugins=grpc,paths=source_relative:./proto/ echo.proto && \
    go build -o /echo-service echo/server

CMD ["/echo-service"]