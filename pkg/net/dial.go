package net

import (
	"context"
	"errors"
	"github.com/Brain-Wave-Ecosystem/go-common/pkg/retry"
	"net"
	"syscall"
	"time"
)

func Dial(network, addr string) (net.Conn, error) {
	dialer := &net.Dialer{}

	conn, err := dialer.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func DialTimeout(network, addr string, timeout, checkInterval time.Duration) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   timeout,
		KeepAlive: checkInterval,
	}

	conn, err := dialer.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	dialer := &net.Dialer{}

	conn, err := dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func DialRetry(ctx context.Context, network, addr string, maxRetry uint, checkInterval time.Duration) (conn net.Conn, err error) {
	dialer := &net.Dialer{
		KeepAlive: checkInterval,
	}

	withRetry := retry.NewRetry(retry.StandardConfig).
		WithAttempts(maxRetry).
		WithRetryIf(func(err error) bool {
			return errors.Is(err, syscall.ECONNREFUSED)
		})

	err = withRetry.Do(func() error {
		conn, err = dialer.DialContext(ctx, network, addr)
		return err
	})

	return conn, nil
}
