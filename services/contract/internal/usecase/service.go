package usecase

import (
	"context"
	"errors"
	"microservices/pkg/validator"
	"microservices/services/contract/internal/domain"
	"microservices/services/contract/internal/repository"
)

var (
	ErrFailedValidation = errors.New("validation failed")
	ErrDuplicate        = errors.New("record duplication")
)

type CreateContractDTO struct {
	Title string `json:"title"`
	Desc  string `json:"description"`
}

type ContractService interface {
	CreateContract(ctx context.Context, input CreateContractDTO) error
	GetContractByID(ctx context.Context, id int64) (*domain.Contract, error)
	GetContracts(ctx context.Context, title string, filters repository.Filters) ([]*domain.Contract, error)
}

type service struct {
	repo repository.Contract
}

func New(repo repository.Contract) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateContract(ctx context.Context, input CreateContractDTO) error {
	contract := domain.Contract{
		Title: input.Title,
		Desc:  input.Desc,
	}

	v := validator.New()

	if ValidateBook(v, &contract); !v.Valid() {
		return ErrFailedValidation
	}

	err := s.repo.Create(ctx, &contract)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetContractByID(ctx context.Context, id int64) (*domain.Contract, error) {
	contract, err := s.repo.GetByID(ctx, id)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return nil, repository.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return contract, nil
}

func (s *service) GetContracts(ctx context.Context, title string, filters repository.Filters) ([]*domain.Contract, error) {
	v := validator.New()

	if repository.ValidateFilters(v, filters); !v.Valid() {
		return nil, ErrFailedValidation
	}
	var contracts []*domain.Contract

	contracts, err := s.repo.GetAll(ctx, title, filters)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return nil, repository.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return contracts, err
}

func ValidateBook(v *validator.Validator, contract *domain.Contract) {
	v.Check(contract.Title != "", "title", "must be provided")
	v.Check(len(contract.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(contract.Desc != "", "description", "must be provided")
	v.Check(len(contract.Desc) >= 1500, "description", "must be greater than 1500 characters")
}
