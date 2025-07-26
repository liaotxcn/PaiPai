goctl model mongo --type chatLog --dir ./apps/im/immodels

goctl rpc protoc apps/im/rpc/im.proto --go_out=./apps/im/rpc --go-grpc_out=./apps/im/rpc --zrpc_out=./apps/im/rpc