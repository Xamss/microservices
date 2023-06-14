package grpc

import (
	contract "microservices/services/contract/internal/delivery/grpc/interface"
	"microservices/services/contract/internal/usecase"
)

type Delivery struct {
	contract.UnimplementedContractServiceServer
	ucContact usecase.ContractService

	options Options
}

type Options struct{}

func New(ucContact usecase.ContractService, o Options) *Delivery {
	var d = &Delivery{
		ucContact: ucContact,
	}

	d.SetOptions(o)
	return d
}

func (d *Delivery) SetOptions(options Options) {
	if d.options != options {
		d.options = options
	}
}
