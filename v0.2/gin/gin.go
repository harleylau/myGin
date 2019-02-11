package gin

import (
	"github.com/julienschmidt/httprouter"

	"net/http"
	"path"
)

type (
	// HandlerFunc .
	HandlerFunc func(*Context)

	// Context .
	Context struct {
		Req      *http.Request
		Writer   http.ResponseWriter
		Params   httprouter.Params
		handlers []HandlerFunc
		engine   *Engine
		index    int8
	}

	// RouterGroup .
	RouterGroup struct {
		Handlers []HandlerFunc
		prefix   string
		parent   *RouterGroup
		engine   *Engine
	}

	// Engine .
	Engine struct {
		*RouterGroup
		router *httprouter.Router
	}
)

// New Engine
func New() *Engine {
	engine := &Engine{}
	engine.RouterGroup = &RouterGroup{nil, "", nil, engine}
	engine.router = httprouter.New()
	return engine
}

// ServeHTTP makes the router implement the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	engine.router.ServeHTTP(w, req)
}

// Run .
func (engine *Engine) Run(addr string) {
	http.ListenAndServe(addr, engine)
}

/************************************/
/********** ROUTES GROUPING *********/
/************************************/

func (group *RouterGroup) createContext(w http.ResponseWriter, req *http.Request, params httprouter.Params, handlers []HandlerFunc) *Context {
	return &Context{
		Writer:   w,
		Req:      req,
		index:    -1,
		engine:   group.engine,
		Params:   params,
		handlers: handlers,
	}
}

// Use .
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.Handlers = append(group.Handlers, middlewares...)
}

// Group .
func (group *RouterGroup) Group(component string, handlers ...HandlerFunc) *RouterGroup {
	prefix := path.Join(group.prefix, component)
	return &RouterGroup{
		Handlers: group.combineHandlers(handlers),
		parent:   group,
		prefix:   prefix,
		engine:   group.engine,
	}
}

// Handle .
func (group *RouterGroup) Handle(method, p string, handlers []HandlerFunc) {
	p = path.Join(group.prefix, p)
	handlers = group.combineHandlers(handlers)
	group.engine.router.Handle(method, p, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		group.createContext(w, req, params, handlers).Next()
	})
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (group *RouterGroup) POST(path string, handlers ...HandlerFunc) {
	group.Handle("POST", path, handlers)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (group *RouterGroup) GET(path string, handlers ...HandlerFunc) {
	group.Handle("GET", path, handlers)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (group *RouterGroup) DELETE(path string, handlers ...HandlerFunc) {
	group.Handle("DELETE", path, handlers)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (group *RouterGroup) PATCH(path string, handlers ...HandlerFunc) {
	group.Handle("PATCH", path, handlers)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (group *RouterGroup) PUT(path string, handlers ...HandlerFunc) {
	group.Handle("PUT", path, handlers)
}

func (group *RouterGroup) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	s := len(group.Handlers) + len(handlers)
	h := make([]HandlerFunc, 0, s)
	h = append(h, group.Handlers...)
	h = append(h, handlers...)
	return h
}

// Next .
func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Writes the given string into the response body and sets the Content-Type to "text/plain"
func (c *Context) String(code int, msg string) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(msg))
}
