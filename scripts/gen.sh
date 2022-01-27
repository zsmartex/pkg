protoc -I=proto --go_out=./gen proto/utils.proto

protoc --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=$PWD $PWD/proto/engine.proto
protoc --go_out=$PWD $PWD/proto/engine.proto
protoc -I=$PWD --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=$GOPATH/src/ $PWD/proto/*.proto
protoc -I=$PWD --go-grpc_out=$GOPATH/src/ $PWD/proto/*.proto
protoc -I=$PWD --go-grpc_out=$GOPATH/src/ $PWD/proto/*.proto