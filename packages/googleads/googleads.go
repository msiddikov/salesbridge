package googleads

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type (
	Service struct {
		ClientID     string
		ClientSecret string
		RedirectURL  string

		Scope                  string
		SaveConnection         func(Connection) (Connection, error)
		UpdateConnectionTokens func(Connection) error
		GetConnectionByID      func(uint) (Connection, error)
		now                    func() time.Time
		tokenExpiryLeeway      time.Duration
	}

	Client struct {
		CustomerInfo CustomerInfo

		ctx      context.Context
		grpcConn *grpc.ClientConn
	}

	CustomerInfo struct {
		CustomerID       string
		ClientCustomerID string
		ManagerID        string
	}
)

func normalizeCustomerID(id string) string {
	if id == "" {
		return ""
	}
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, id)
}

func NewService(
	ClientId, ClientSecret, RedirectURL string,
	saveConnection func(Connection) (Connection, error),
	updateConnectionTokens func(Connection) error,
	getConnectionByID func(uint) (Connection, error),
) Service {
	svc := Service{
		ClientID:     ClientId,
		ClientSecret: ClientSecret,
		RedirectURL:  RedirectURL,

		SaveConnection:         saveConnection,
		UpdateConnectionTokens: updateConnectionTokens,
		GetConnectionByID:      getConnectionByID,
		now:                    time.Now,
		tokenExpiryLeeway:      1 * time.Minute,
	}
	return svc
}

func (s *Service) NewClient(connId uint, customerInfo CustomerInfo) (*Client, error) {
	customerInfo.CustomerID = normalizeCustomerID(customerInfo.CustomerID)
	customerInfo.ClientCustomerID = normalizeCustomerID(customerInfo.ClientCustomerID)
	customerInfo.ManagerID = normalizeCustomerID(customerInfo.ManagerID)

	conn, err := s.GetConnectionByID(connId)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	ctx, _, err = s.WithHeaders(ctx, conn)
	if err != nil {
		return nil, err
	}

	grpcConn, err := grpc.NewClient("googleads.googleapis.com:443", grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		return nil, err
	}

	customerID := customerInfo.CustomerID
	if customerID == "" {
		customerID = customerInfo.ClientCustomerID
	}
	if customerID == "" {
		return nil, fmt.Errorf("googleads: customer id required")
	}
	customerInfo.CustomerID = customerID

	if customerInfo.ManagerID != "" {
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			mdCopy := md.Copy()
			mdCopy.Set("login-customer-id", customerInfo.ManagerID)
			ctx = metadata.NewOutgoingContext(ctx, mdCopy)
		} else {
			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("login-customer-id", customerInfo.ManagerID))
		}
	}

	return &Client{
		CustomerInfo: customerInfo,
		grpcConn:     grpcConn,
		ctx:          ctx,
	}, nil
}

func (c *Client) Close() error {
	return c.grpcConn.Close()
}
