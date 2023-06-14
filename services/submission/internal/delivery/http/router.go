package http

import (
	"microservices-go/services/order/internal/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type router struct {
	order OrderHandler
}

func NewRouter(orderService service.OrderService) *router {
	return &router{order: *NewHandler(orderService)}
}

func (r *router) GetRoutes() http.Handler {

	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/v1/submission/create", r.order.CreateOrder)
	router.HandlerFunc(http.MethodPost, "/v1/submission/show", r.order.ShowOrder)

	return router
}
