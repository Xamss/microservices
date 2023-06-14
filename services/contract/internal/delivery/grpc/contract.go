package grpc

import (
	"context"

	contact "microservices/services/contract/internal/delivery/grpc/interface"
)

func (d *Delivery) CreateGroup(ctx context.Context, request *contact.CreateContractRequest) (*contact.CreateContractResponse, error) {
	panic("implement me")
}

func (d *Delivery) UpdateGroup(ctx context.Context, request *contact.UpdateContractRequest) (*contact.UpdateContractResponse, error) {
	panic("implement me")
}

func (d *Delivery) DeleteGroup(ctx context.Context, request *contact.DeleteContractRequest) (*contact.DeleteContractResponse, error) {
	panic("implement me")
}
