goctl api go -api apps/eino_chat/api/eino.api -dir apps/eino_chat/api -style gozero

goctl rpc protoc apps/eino_chat/rpc/eino.proto --go_out=./apps/eino_chat/rpc --go-grpc_out=./apps/eino_chat/rpc --zrpc_out=./apps/eino_chat/rpc