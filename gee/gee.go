package gee

import (
	"net/http"
)

type HandlerFunc func(*Context)

type Engine struct {
	router *router
}

func New() *Engine  {
	return &Engine{
		router: newRouter(),
	}
}

func (engine *Engine)addRoute(method,pattern string,handler HandlerFunc)  {
	key := method + "-" + pattern
	engine.router.handlers[key] = handler
}


func (engine *Engine)GET(pattern string, handler HandlerFunc)  {
	engine.router.addRoute("GET",pattern,handler)
}

func (engine *Engine)POST(pattern string, handler HandlerFunc)  {
	engine.router.addRoute("POST",pattern,handler)
}

func (engine *Engine)Run(addr string) (err error)  {
	return http.ListenAndServe(addr,engine)
}

func (engine *Engine)ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	engine.router.handle(newContext(w,req))
}

