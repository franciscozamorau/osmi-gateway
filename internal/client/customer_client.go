// internal/client/customer_client.go
package client

import (
	"context"

	pb "github.com/franciscozamorau/osmi-protobuf/gen/pb"
	"google.golang.org/grpc"
)

// CustomerClient encapsula el cliente gRPC para customers
type CustomerClient struct {
	client pb.OsmiServiceClient
}

// NewCustomerClient crea un nuevo cliente de customers
func NewCustomerClient(conn *grpc.ClientConn) *CustomerClient {
	return &CustomerClient{
		client: pb.NewOsmiServiceClient(conn),
	}
}

// CreateCustomer llama al método gRPC CreateCustomer
func (c *CustomerClient) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CustomerResponse, error) {
	return c.client.CreateCustomer(ctx, req)
}

// GetCustomer llama al método gRPC GetCustomer
func (c *CustomerClient) GetCustomer(ctx context.Context, req *pb.CustomerLookup) (*pb.CustomerResponse, error) {
	return c.client.GetCustomer(ctx, req)
}

// UpdateCustomer llama al método gRPC UpdateCustomer
func (c *CustomerClient) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (*pb.CustomerResponse, error) {
	return c.client.UpdateCustomer(ctx, req)
}

// ListCustomers llama al método gRPC ListCustomers
func (c *CustomerClient) ListCustomers(ctx context.Context, req *pb.ListCustomersRequest) (*pb.CustomerListResponse, error) {
	return c.client.ListCustomers(ctx, req)
}

// GetCustomerStats llama al método gRPC GetCustomerStats
func (c *CustomerClient) GetCustomerStats(ctx context.Context, req *pb.Empty) (*pb.CustomerStatsResponse, error) {
	return c.client.GetCustomerStats(ctx, req)
}
