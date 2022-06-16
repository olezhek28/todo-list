package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/olezhek28/todo-list/middleware"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/task/v1/list", middleware.GetList).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/task/v1/create", middleware.Create).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/task/v1/done/{id}", middleware.Done).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/task/v1/undone/{id}", middleware.Undone).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/task/v1/delete/{id}", middleware.Delete).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/task/v1/delete-all", middleware.DeleteAll).Methods(http.MethodDelete, http.MethodOptions)

	return router
}
