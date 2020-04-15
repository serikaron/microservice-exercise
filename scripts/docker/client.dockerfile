FROM serika-registry:443/go-dev as development

COPY pkg /mse/pkg
COPY proto/*.proto proto/gen.go /mse/proto/
COPY client /mse/client

RUN go generate mse/proto && \
    go build -o /client mse/client

#FROM alpine:3.11 as release
#COPY res /mse/res
#COPY --from=development /client /
#ENTRYPOINT /client --auth-host "" --auth-port $auth_port --cert-path $cert_path --chat-host "" --chat-port $chat_port

