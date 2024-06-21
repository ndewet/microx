package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseWriter interface {
	Write([]byte) (int, error)
	WriteHeader(int)
	Header() http.Header
}

type Response interface {
	Write(ResponseWriter) error
}

type RawResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

func (response RawResponse) Write(writer ResponseWriter) error {
	writer.WriteHeader(response.StatusCode)
	for key, value := range response.Headers {
		writer.Header().Set(key, value)
	}
	_, err := writer.Write(response.Body)
	return err
}

type ObjectResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       interface{}
}

func (response ObjectResponse) Write(writer ResponseWriter) error {
	serializedBody, err := json.Marshal(response.Body)
	if err != nil {
		return InternalServerError{err}.Write(writer)
	}
	return RawResponse{
		StatusCode: response.StatusCode,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       serializedBody,
	}.Write(writer)
}

type JSONResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       map[string]interface{}
}

func (response JSONResponse) Write(writer ResponseWriter) error {
	return ObjectResponse{
		StatusCode: response.StatusCode,
		Headers:    response.Headers,
		Body:       response.Body,
	}.Write(writer)
}

type ErrorResponse struct {
	StatusCode int
	Message    string
	Error      error
}

func (response ErrorResponse) Write(writer ResponseWriter) error {
	var body []byte
	if response.Error == nil {
		body = []byte(response.Message)
	} else {
		body = []byte(fmt.Sprintf("%s: %s", response.Message, response.Error))
	}
	return RawResponse{
		StatusCode: response.StatusCode,
		Body:       body,
	}.Write(writer)
}

type InternalServerError struct {
	Error error
}

func (response InternalServerError) Write(writer ResponseWriter) error {
	return ErrorResponse{
		StatusCode: 500,
		Message:    "internal server error",
		Error:      response.Error,
	}.Write(writer)
}

type ServiceUnavailable struct{}

func (response ServiceUnavailable) Write(writer ResponseWriter) error {
	return RawResponse{
		StatusCode: 503,
		Body:       []byte("service unavailable"),
	}.Write(writer)
}

type BadRequest struct {
	Error error
}

func (response BadRequest) Write(writer ResponseWriter) error {
	return ErrorResponse{
		StatusCode: 400,
		Message:    "bad request",
		Error:      response.Error,
	}.Write(writer)
}
