package http

import (
	"microservices/services/contract/internal/usecase"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type router struct {
	contract ContractHandler
}

func NewRouter(bookService usecase.ContractService) *router {
	return &router{contract: *NewHandler(bookService)}
}

func (r *router) GetRoutes() http.Handler {

	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/v1/books", r.contract.CreateContractHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books/:id", r.contract.ShowContractHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books", r.contract.ListContractHandler)

	return router
}
