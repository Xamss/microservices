package http

import (
	"errors"
	"microservices/pkg/request"
	"microservices/pkg/validator"
	"microservices/services/contract/internal/repository"
	"microservices/services/contract/internal/usecase"
	"net/http"
)

type ContractHandler struct {
	contractService usecase.ContractService
}

func NewHandler(service usecase.ContractService) *ContractHandler {
	return &ContractHandler{contractService: service}
}

func (h *ContractHandler) CreateContractHandler(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateContractDTO

	err := request.ReadJSON(w, r, &input)
	if err != nil {
		request.BadRequestResponse(w, r, err)
		return
	}

	contract := usecase.CreateContractDTO{
		Title: input.Title,
		Desc:  input.Desc,
	}

	err = h.contractService.CreateContract(r.Context(), contract)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrFailedValidation):
			request.BadRequestResponse(w, r, err)
			return
		case errors.Is(err, usecase.ErrDuplicate):
			request.RecordDuplicationResponse(w, r)
			return
		default:
			request.ServerErrorResponse(w, r, err)
			return
		}
	}

	err = request.WriteJSON(w, http.StatusCreated, map[string]any{"contract": contract}, nil)
	if err != nil {
		request.ServerErrorResponse(w, r, err)
		return
	}

}

func (h *ContractHandler) ShowContractHandler(w http.ResponseWriter, r *http.Request) {
	id, err := request.ReadIDParam(r)
	if err != nil {
		request.NotFoundResponse(w, r)
		return
	}

	contract, err := h.contractService.GetContractByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			request.NotFoundResponse(w, r)
			return
		default:
			request.ServerErrorResponse(w, r, err)
			return
		}
	}
	err = request.WriteJSON(w, http.StatusOK, map[string]any{"contract": contract}, nil)
	if err != nil {
		request.ServerErrorResponse(w, r, err)
		return
	}
}

func (h *ContractHandler) ListContractHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string
		repository.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	input.Title = request.ReadString(qs, "title", "")

	input.Filters.Page = request.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = request.ReadInt(qs, "page_size", 20, v)
	input.Filters.Sort = request.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "-id", "-title"}

	contracts, err := h.contractService.GetContracts(r.Context(), input.Title, input.Filters)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			request.NotFoundResponse(w, r)
			return
		default:
			request.ServerErrorResponse(w, r, err)
			return
		}
	}
	err = request.WriteJSON(w, http.StatusOK, map[string]any{"contracts": contracts}, nil)
	if err != nil {
		request.ServerErrorResponse(w, r, err)
		return
	}
}
