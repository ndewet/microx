package httpx

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func MockHandler(Request) (Response, error) {
	return &MockResponse{}, nil
}

type MockHandlerDetails struct {
	HandlerCallCount   int
	ResponseWriteCount int
}

func CreateMockHandler() (Handler, *MockHandlerDetails) {
	details := MockHandlerDetails{
		HandlerCallCount:   0,
		ResponseWriteCount: 0,
	}
	return func(r Request) (Response, error) {
		details.HandlerCallCount++
		return &MockResponse{
			WriteCounter: &details.ResponseWriteCount,
		}, nil
	}, &details
}

func CreateErrorMockHandler() (Handler, *MockHandlerDetails) {
	details := MockHandlerDetails{
		HandlerCallCount:   0,
		ResponseWriteCount: 0,
	}
	return func(r Request) (Response, error) {
		details.HandlerCallCount++
		return nil, fmt.Errorf("error")
	}, &details
}

func CreatePanickingMockHandler() (Handler, *MockHandlerDetails) {
	details := MockHandlerDetails{
		HandlerCallCount:   0,
		ResponseWriteCount: 0,
	}
	return func(r Request) (Response, error) {
		details.HandlerCallCount++
		panic("panic")
	}, &details
}

func TestAdaptedHandlerCatchesPanicsAndReturnsInternalServerError(t *testing.T) {
	handler, _ := CreatePanickingMockHandler()
	adaptedHandler := adapt(handler)
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	adaptedHandler.ServeHTTP(writer, request)
	if writer.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 500, got %d", writer.Code)
	}
}

func TestAdaptedHandlerCatchesErrorsAndReturnsInternalServerError(t *testing.T) {
	handler, _ := CreateErrorMockHandler()
	adaptedHandler := adapt(handler)
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	adaptedHandler.ServeHTTP(writer, request)
	if writer.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 500, got %d", writer.Code)
	}
}
