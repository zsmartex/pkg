GOPATH=~/go

protoc -I=$PWD --go_out=$GOPATH/src/ $PWD/proto/utils.proto
protoc -I=$PWD --go_out=$GOPATH/src/ $PWD/proto/symbol.proto
protoc -I=$PWD --go_out=$GOPATH/src/ $PWD/proto/order.proto
protoc -I=$PWD --go_out=$GOPATH/src/ $PWD/proto/engine.proto
protoc -I=$PWD --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=$GOPATH/src/ $PWD/proto/engine.proto
protoc -I=$PWD --go_out=$GOPATH/src/ $PWD/proto/quantex.proto
protoc -I=$PWD --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=$GOPATH/src/ $PWD/proto/quantex.proto
