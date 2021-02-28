package pb

//go:generate protoc -I=../../proto --go_out=. --go_opt=module=github.com/bbengfort/cosmos/pkg/pb --go-grpc_out=. --go-grpc_opt=module=github.com/bbengfort/cosmos/pkg/pb cosmos/v1alpha1/api.proto
