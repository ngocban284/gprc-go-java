package client

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptorClient struct {
	authClient  *AuthClient
	authMethod  map[string]bool
	accessToken string
}

func NewAuthInterceptorClient(
	authClient *AuthClient,
	authMethod map[string]bool,
	refreshDuration time.Duration,
) (*AuthInterceptorClient, error) {
	interceptor := &AuthInterceptorClient{
		authClient: authClient,
		authMethod: authMethod,
	}

	err := interceptor.scheduleRefreshToken(refreshDuration)
	if err != nil {
		return nil, err
	}

	return interceptor, nil
}

func (c *AuthInterceptorClient) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", c.accessToken)
}

func (c *AuthInterceptorClient) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Print("---> intercepting unary method: ", method)

		if c.authMethod[method] {
			return invoker(c.attachToken(ctx), method, req, reply, cc, opts...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (c *AuthInterceptorClient) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log.Print("---> intercepting stream method: ", method)

		if c.authMethod[method] {
			return streamer(c.attachToken(ctx), desc, cc, method, opts...)
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (c *AuthInterceptorClient) refreshToken() error {
	accessToken, err := c.authClient.Login()
	if err != nil {
		return err
	}

	c.accessToken = accessToken

	return nil
}

func (c *AuthInterceptorClient) scheduleRefreshToken(refreshDuration time.Duration) error {
	err := c.refreshToken()
	if err != nil {
		return err
	}

	go func() {
		for {
			time.Sleep(refreshDuration)
			err := c.refreshToken()
			if err != nil {
				return
			}
		}
	}()

	return nil
}
