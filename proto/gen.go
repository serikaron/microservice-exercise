package proto

//go:generate protoc -I./ --go_out=plugins=grpc,paths=source_relative:. auth.proto
//go:generate protoc -I./ --go_out=plugins=grpc,paths=source_relative:. auth.internal.proto
//go:generate protoc -I./ --go_out=plugins=grpc,paths=source_relative:. chat.proto
