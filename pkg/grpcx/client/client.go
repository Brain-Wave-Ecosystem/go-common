package client

import (
	"fmt"
	"github.com/Brain-Wave-Ecosystem/go-common/pkg/retry"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewServiceClient[T any](addr string, creationFun func(cc grpc.ClientConnInterface) T, opts ...grpc.DialOption) (T, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("dynamic:///%s", addr), opts...)
	if err != nil {
		fmt.Printf("grpc.NewClient failed: %v\n", err)
		return *new(T), nil
	}

	return creationFun(conn), nil
}

func NewServiceClientWithRetry[T any](addr string, creationFun func(cc grpc.ClientConnInterface) T, retryConf *retry.Config, opts ...grpc.DialOption) (T, error) {
	var conn *grpc.ClientConn
	var err error

	withRetry := retry.NewRetry(retryConf).WithRetryIf(func(err error) bool {
		if err != nil {
			grpcErr, found := status.FromError(err)
			if !found {
				grpcErr = status.Convert(err)
			}

			grpcCode := grpcErr.Code()
			if grpcCode == codes.Unavailable {
				return true
			}
		}

		return false
	})

	err = withRetry.Do(func() error {
		conn, err = grpc.NewClient(fmt.Sprintf("dynamic:///%s", addr), opts...)
		return err
	})
	if err != nil {
		fmt.Printf("grpc.NewClient creation failed: %v\n", err)
		return *new(T), nil
	}

	return creationFun(conn), nil
}
