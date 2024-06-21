package httpx

import (
	"context"
	"net/http"
)

// Server is an HTTP server that wraps an http.Server and a Router.
// The server is not started automatically.
type Server struct {
	// server is the underlying HTTP server.
	server *http.Server
	// middleware is the list of middleware applied to the server.
	middleware []Middleware
	// router is the router used by the server.
	router *Router
}

// NewServer creates a new HTTP server listening on the given address.
// The server is not started automatically.
func NewServer(address string) *Server {
	router := NewRouter()
	middleware := []Middleware{}
	httpServer := &http.Server{Addr: address}
	server := &Server{server: httpServer, middleware: middleware, router: router}
	return server
}

// WithRouter sets the router for the server.
// Overwrites any existing router.
func (server *Server) WithRouter(router *Router) *Server {
	server.router = router
	return server
}

// WithMiddleware adds a middleware to the server.
func (server *Server) WithMiddleware(middleware Middleware) *Server {
	// Prepend middleware to ensure that the middleware is executed in the correct order.
	// The last applied middleware is the first to be executed.
	server.middleware = append([]Middleware{middleware}, server.middleware...)
	return server
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
// It blocks until all connections are closed.
//
// see http.Server.Shutdown for more details.
func (server *Server) Shutdown() error {
	return server.server.Shutdown(context.Background())
}

// ForceShutdown forcefully shuts down the server.
// It does not wait for connections to close.
//
// see http.Server.Close for more details.
func (server *Server) ForceShutdown() {
	server.server.Close()
}

// Start starts the server.
// It blocks until the server is shut down.
//
// see http.Server.ListenAndServe for more details.
func (server *Server) Start() error {
	server.server.Handler = server.router
	for _, middleware := range server.middleware {
		server.server.Handler = middleware(server.server.Handler)
	}
	return server.server.ListenAndServe()
}
