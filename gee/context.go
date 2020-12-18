package gee

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

type Context struct {
	Writer http.ResponseWriter
	Req *http.Request
	Params map[string]string
	Path string
	Method string
	StatusCode int
	
	Handlers []HandlerFunc
	index int
}
type H map[string]interface{}
func newContext(w http.ResponseWriter,req *http.Request) *Context{
	return &Context{
		Writer: w,
		Req: req,
		Path: req.URL.Path,
		Method: req.Method,
		Handlers: make([]HandlerFunc,0),
		index: -1,
	}
}

func (c *Context)PostForm(key string) string  {
	return c.Req.FormValue(key)
}

func (c *Context)Query(key string) string  {
	return c.Req.URL.Query().Get(key)
}

func (c *Context)Status(code int)   {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context)String(code int,format string,values ...interface{})  {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format,values...)))
}

func (c *Context)JSON(code int,obj interface{})  {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(),500)
	}
}
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) Next() {
	c.index++
	for ;c.index < len(c.Handlers);c.index++{
		c.Handlers[c.index](c)
	}
}

func (c *Context) Abort() {
	c.index = math.MaxInt8
}

