package httpx

import (
	"fmt"
	"net/http"
)

// Response is a type alias for http.Response.
type Request http.Request

// Handler is any function that handles an HTTP request and returns an (Response, error) tuple.
type Handler = func(Request) (Response, error)

type Router struct {
	multiplexer *http.ServeMux
}

func NewRouter() *Router {
	multiplexer := http.NewServeMux()
	return &Router{multiplexer}
}

// CreateRouter creates a router from a map of routes and a map of routers.
// The routes map maps methods to paths to handlers.
// The routers map maps paths to routers.
// The paths must not overlap.
func CreateRouter(routes map[Method]map[string]Handler, routers map[string]*Router) *Router {
	router := NewRouter()
	for method, routes := range routes {
		for path, handler := range routes {
			router.Route(method, path, handler)
		}
	}
	for path, subRouter := range routers {
		router.Link(path, subRouter)
	}
	return router
}

// CreateRouterWithPrefix creates a router with the given prefix from a map of routes and a map of routers.
//
// The prefix must be a valid path.
// Useful for creating versioned APIs.
func CreateRouterWithPrefix(prefix string, routes map[Method]map[string]Handler, routers map[string]*Router) *Router {
	return NewRouter().Link(
		prefix,
		CreateRouter(routes, routers),
	)

}

// Route registers a handler for the given method and path.
// The path must start with a "/" and end with a "/".
// The path must not contain spaces.
// The path must not contain consecutive slashes.
func (router *Router) Route(method Method, path string, handler Handler) *Router {
	validate(path)
	pattern := fmt.Sprintf("%s %s", method, path)
	router.multiplexer.Handle(pattern, adapt(handler))
	return router
}

// Link links the otherRouter router to the path.
// The path must not already be handled by the router.
// When path == "/", this is equivalent to merging the routers.
func (router *Router) Link(path string, otherRouter *Router) *Router {
	if path == "/" {
		router.multiplexer.Handle(path, otherRouter.multiplexer)
		return router
	}
	validate(path)
	prefix := path[:len(path)-1]
	router.multiplexer.Handle(path, http.StripPrefix(prefix, otherRouter.multiplexer))
	return router
}

func (router *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	router.multiplexer.ServeHTTP(writer, request)
}

// adapt adapts a handler to the http.HandlerFunc interface by ensuring the the return response is written to the http.ResponseWriter.
// It recovers from panics and any errors returned by the handler.
// Both cases result in an 500 response.
func adapt(handler Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				InternalServerError{fmt.Errorf("internal server error")}.Write(writer)
			}
		}()
		response, err := handler(Request(*request))
		if err != nil {
			response = InternalServerError{fmt.Errorf("internal server error")}
		}
		response.Write(writer)
	})
}

// validate validates the path.
// The path must start with a "/" and end with a "/".
// The path must not contain spaces.
// The path must not contain consecutive slashes.
// The path must not be empty.
// Panics if the path is invalid.
func validate(path string) {
	if len(path) == 0 {
		panic("path can not be empty")
	}
	if path[0] != '/' {
		panic("path must start with /")
	}
	if path[len(path)-1] != '/' {
		panic("path must end with /")
	}
	for i := 1; i < len(path)-1; i++ {
		if path[i] == '/' && path[i+1] == '/' {
			panic("path must not contain consecutive slashes")
		}
	}
	for i := 0; i < len(path); i++ {
		if path[i] == ' ' {
			panic("path must not contain spaces")
		}
	}
}
