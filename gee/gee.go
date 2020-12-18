package gee

import (
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type Engine struct {
	*RouterGroup
	groups []*RouterGroup
	router *router
}

type RouterGroup struct {
	prefix string
	engine *Engine
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
}

func New() *Engine  {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine:engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
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
	var midwares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path,group.prefix) {
			midwares = append(midwares,group.middlewares...)
		}
	}
	c := newContext(w,req)
	c.Handlers = midwares
	engine.router.handle(c)
}

func (group *RouterGroup)Group(prefix string) *RouterGroup  {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups,newGroup)
	return newGroup
}

func (group *RouterGroup)Use(f ...HandlerFunc)  {
	group.middlewares = append(group.middlewares,f...)
}

func (group *RouterGroup)addRoute(method,pattern string,handler HandlerFunc)  {
	finalPattern := group.prefix + pattern
	group.engine.router.addRoute(method,finalPattern,handler)
}


func (group *RouterGroup)GET(pattern string, handler HandlerFunc)  {
	group.addRoute("GET",pattern,handler)
}

func (group *RouterGroup)POST(pattern string, handler HandlerFunc)  {
	group.addRoute("POST",pattern,handler)
}

