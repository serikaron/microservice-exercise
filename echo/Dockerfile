FROM grpc-go:latest as development

COPY . /echo/

WORKDIR /echo

RUN go test echo/server

RUN mkdir -p /echo/proto && \
#    go generate && \
#    protoc -I../proto --go_out=plugins=grpc,paths=source_relative:../proto/ echo.proto && \
    protoc -I./proto --go_out=plugins=grpc,paths=source_relative:./proto/ echo.proto && \
    go build -o /echo-service echo/server && \
    go build -o /echo-client echo/client

FROM alpine:3.11 as echo-service
COPY --from=development /echo-service /
EXPOSE 55555
CMD ["/echo-service"]

FROM alpine:3.11 as echo-client
COPY --from=development /echo-client /
EXPOSE 55555
CMD ["/echo-client"]
