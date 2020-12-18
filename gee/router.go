package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router  {
	return &router{
		handlers: make(map[string]HandlerFunc),
		roots: make(map[string]*node),
	}
}
func parsePattern(pattern string)[]string  {
	vs := strings.Split(pattern,"/")
	parts := make([]string,0)
	for _,part := range vs{
		if part != ""{
			 parts = append(parts,part)
			if part[0] == '*' {
				break
			}
		}
	}
	return parts
}
func (r *router)addRoute(method,pattern string,handler HandlerFunc)  {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_,ok := r.roots[method]
	if  !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern,parts,0)
	r.handlers[key] = handler
}
//例如/p/go/doc匹配到/p/:lang/doc，解析结果为：{lang: "go"}，/static/css/geektutu.css匹配到/static/*filepath，解析结果为{filepath: "css/geektutu.css"}。
func (r *router)getRoute(method,path string) (*node,map[string]string) {
	searchParts := parsePattern(path)
	parmes := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil,nil
	}
	n := root.search(searchParts,0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for i, part := range parts {
			if part[0] == ':' {
				parmes[part[1:]] = searchParts[i]
			}
			if part[0] == '*' && len(part) > 1 {
				parmes[part[1:]] = strings.Join(searchParts[i:],"/")
				break
			}
		}
		return n,parmes
	}
	return nil, nil
}

func (r *router)handle(c *Context)  {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.Handlers = append(c.Handlers,r.handlers[key])
	} else {
		c.Handlers = append(c.Handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}