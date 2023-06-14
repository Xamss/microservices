package http

import (
	"errors"
	"fmt"
	"microservices/pkg/request"
	"microservices/services/submission/internal/repository"
	"microservices/services/submission/internal/usecase"
	"net/http"
)

type OrderHandler struct {
	orderService usecase.OrderService
}

func NewHandler(service usecase.OrderService) *OrderHandler {
	return &OrderHandler{orderService: service}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var dto usecase.CreateOrderDTO

	if err := request.ReadJSON(w, r, &dto); err != nil {
		request.BadRequestResponse(w, r, err)
		return
	}

	input := usecase.CreateOrderDTO{
		BookID: dto.BookID,
		Email:  dto.Email,
	}

	err := h.orderService.Create(r.Context(), input)
	if err != nil {
		request.ServerErrorResponse(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (h *OrderHandler) ShowOrder(w http.ResponseWriter, r *http.Request) {
	//email, err := request.ReadEmailParam(r)
	//if err != nil {
	//	request.ServerErrorResponse(w, r, err)
	//	return
	//}

	var dto usecase.CreateOrderDTO
	fmt.Println("ShowOrder")
	if err := request.ReadJSON(w, r, &dto); err != nil {
		request.BadRequestResponse(w, r, err)
		fmt.Println(dto)
		return
	}
	fmt.Println(dto)
	orders, err := h.orderService.Show(r.Context(), dto.Email)
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
	err = request.WriteJSON(w, http.StatusOK, map[string]any{"orders": orders}, nil)
	if err != nil {
		request.ServerErrorResponse(w, r, err)
		return
	}
}
