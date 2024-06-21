package httpx

import (
	"fmt"
	"net/http"
	"regexp"
)

type Multiplexer interface {
	Handle(pattern string, handler http.Handler)
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}

type Router struct {
	multiplexer Multiplexer
}

func NewRouter() *Router {
	multiplexer := Multiplexer(http.NewServeMux())
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

func validate(path string) {
	regex := `^\/(?:[^\/\s{}]+|{[^\/\s{}]+})*(?:\/(?:[^\/\s{}]+|{[^\/\s{}]+}))*\/$`
	if !regexp.MustCompile(regex).MatchString(path) {
		panic("path is invalid")
	}
}
