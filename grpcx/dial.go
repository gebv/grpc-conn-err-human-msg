package grpcx

import (
	"context"
	"crypto/tls"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Dial(parent context.Context, addr string, timeout time.Duration, sets ...DialOption) (*grpc.ClientConn, context.CancelFunc, error) {
	if err := TCPConnectOK(timeout, addr); err != nil {
		return nil, nil, err
	}

	opts := DialOptions{}
	for _, set := range sets {
		set(&opts)
	}

	ctx, cancel := context.WithCancel(parent)

	if !opts.PlainText {
		verifyCert := NewManualVerifyPeerCertificate(opts.SkipTLSVerify, opts.Fingerprint, "")
		opts.GRPCOptions = append(opts.GRPCOptions, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			// it is ok we used VerifyPeerCertificate
			InsecureSkipVerify:    true,
			VerifyPeerCertificate: verifyCert.VerifyPeerCertificateOption(ctx, cancel),
		},
		)))
	}

	conn, err := grpc.DialContext(ctx, addr, opts.GRPCOptions...)
	if err != nil {
		cancel()
		return nil, nil, err
	}
	return conn, cancel, nil
}

type DialOption func(opts *DialOptions)

type DialOptions struct {
	// if true then disabled TLS
	PlainText bool
	// skips certificate verification
	SkipTLSVerify bool
	// if set then checks fingerprint the server cert
	Fingerprint string
	// standart grpc opts
	GRPCOptions []grpc.DialOption
}

func PlainText() DialOption {
	return func(opts *DialOptions) {
		opts.PlainText = true
	}
}

func Fingerprint(in string) DialOption {
	return func(opts *DialOptions) {
		opts.Fingerprint = in
	}
}

func SkipTLSVerify() DialOption {
	return func(opts *DialOptions) {
		opts.SkipTLSVerify = true
	}
}

func AddStdGRPCOptions(std ...grpc.DialOption) DialOption {
	return func(opts *DialOptions) {
		opts.GRPCOptions = std
	}
}
