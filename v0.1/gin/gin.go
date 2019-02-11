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
		Req     *http.Request
		Writer  http.ResponseWriter
		Params  httprouter.Params
		handler HandlerFunc
		engine  *Engine
	}

	// RouterGroup .
	RouterGroup struct {
		Handler HandlerFunc
		prefix  string
		parent  *RouterGroup
		engine  *Engine
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

func (group *RouterGroup) createContext(w http.ResponseWriter, req *http.Request, params httprouter.Params, handler HandlerFunc) *Context {
	return &Context{
		Writer:  w,
		Req:     req,
		engine:  group.engine,
		Params:  params,
		handler: handler,
	}
}

// Group .
func (group *RouterGroup) Group(component string) *RouterGroup {
	prefix := path.Join(group.prefix, component)
	return &RouterGroup{
		Handler: nil,
		parent:  group,
		prefix:  prefix,
		engine:  group.engine,
	}
}

// Handle .
func (group *RouterGroup) Handle(method, p string, handler HandlerFunc) {
	p = path.Join(group.prefix, p)
	group.engine.router.Handle(method, p, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		group.createContext(w, req, params, handler).Next()
	})
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (group *RouterGroup) POST(path string, handler HandlerFunc) {
	group.Handle("POST", path, handler)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (group *RouterGroup) GET(path string, handler HandlerFunc) {
	group.Handle("GET", path, handler)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (group *RouterGroup) DELETE(path string, handler HandlerFunc) {
	group.Handle("DELETE", path, handler)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (group *RouterGroup) PATCH(path string, handler HandlerFunc) {
	group.Handle("PATCH", path, handler)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (group *RouterGroup) PUT(path string, handler HandlerFunc) {
	group.Handle("PUT", path, handler)
}

// Next .
func (c *Context) Next() {
	c.handler(c)
}

// Writes the given string into the response body and sets the Content-Type to "text/plain"
func (c *Context) String(code int, msg string) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(msg))
}
