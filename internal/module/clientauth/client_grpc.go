package clientauth

import (
	"context"
	"time"

	pb "github.com/avrebarra/goggle/internal/module/clientauth/proto/generated"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var _ Client = (*ClientGRPC)(nil)

type ClientGRPC struct {
	conn    *grpc.ClientConn
	client  pb.AuthServiceClient
	timeout time.Duration
}

type ConfigGRPC struct {
	Address string        `validate:"required"`
	Timeout time.Duration `validate:"required"`
}

func NewClientGRPC(cfg ConfigGRPC) (out *ClientGRPC, err error) {
	conn, err := grpc.Dial(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to auth service")
	}

	return &ClientGRPC{
		conn:    conn,
		client:  pb.NewAuthServiceClient(conn),
		timeout: cfg.Timeout,
	}, nil
}

func (c *ClientGRPC) Close() error {
	return c.conn.Close()
}

func (c *ClientGRPC) ValidateSession(ctx context.Context, token string) (*UserInfo, error) {
	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Call the gRPC service
	resp, err := c.client.ValidateSession(timeoutCtx, &pb.ValidateSessionRequest{
		Token: token,
	})
	if err != nil {
		// Handle gRPC status errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, ErrUnauthorized
			case codes.InvalidArgument:
				return nil, ErrInvalidToken
			case codes.Unavailable:
				return nil, ErrServiceDown
			default:
				return nil, errors.Wrap(err, "auth service error")
			}
		}
		return nil, errors.Wrap(err, "failed to validate session")
	}

	// Check if token is valid
	if !resp.Valid {
		return nil, ErrInvalidToken
	}

	// Map response to domain model
	return &UserInfo{
		UserID:   resp.User.UserId,
		Username: resp.User.Username,
		Email:    resp.User.Email,
	}, nil
}
