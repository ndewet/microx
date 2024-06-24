package httpx

import (
	"fmt"
	"net/http"
)

// Handler is any function that handles an HTTP request and returns an (Response, error) tuple.
type Handler = func(Request) (Response, error)

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
