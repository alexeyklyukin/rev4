package routes

import (
	"net/http"

	"github.com/alexeyklyukin/rev4/pkg/controller"
	"github.com/julienschmidt/httprouter"

)

func Routes(ctl *controller.Controller) *httprouter.Router {
	router := httprouter.New()
	router.Handler(http.MethodGet, "/", http.HandlerFunc(ctl.Index))
	router.Handler(http.MethodPut, "/hello/:name", http.HandlerFunc(ctl.RecordBirthday))
	router.Handler(http.MethodGet, "/hello", http.HandlerFunc(ctl.TellBirthday))
	return router
}