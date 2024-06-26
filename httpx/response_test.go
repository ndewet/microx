package httpx

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

const MOCK_BODY = "Hello, World!"
const MOCK_BODY_JSON = `{"#":"Hello, World!"}`
const CONTENT_TYPE_HEADER_KEY = "Content-Type"
const CONTENT_TYPE_JSON = "application/json"
const CONTENT_TYPE_TEXT = "text/plain"
const EXPECTED_STRING_ERROR = "Expected %s, got %s"
const EXPECTED_DIGIT_ERROR = "Expected %d, got %d"

const SERVICE_UNAVAILABLE = "service unavailable"
const INTERNAL_SERVER_ERROR = "internal server error"

type MockResponse struct {
	WriteCounter *int
}

func (m *MockResponse) Write(ResponseWriter) error {
	*m.WriteCounter += 1
	return nil
}

func TestRawResponseWritesBody(t *testing.T) {
	response := RawResponse{
		StatusCode: 200,
		Body:       []byte(MOCK_BODY),
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	bytes := writer.Body.Bytes()
	if string(bytes) != MOCK_BODY {
		t.Errorf("Expected Hello, World!, got %s", string(bytes))
	}
}

func TestRawResponseWritesHeaders(t *testing.T) {
	response := RawResponse{
		StatusCode: 500,
		Headers:    map[string]string{CONTENT_TYPE_HEADER_KEY: CONTENT_TYPE_TEXT},
		Body:       []byte(""),
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Header().Get(CONTENT_TYPE_HEADER_KEY) != CONTENT_TYPE_TEXT {
		t.Errorf(EXPECTED_STRING_ERROR, CONTENT_TYPE_TEXT, writer.Header().Get(CONTENT_TYPE_HEADER_KEY))
	}
}

func TestRawResponseWritesStatusCode(t *testing.T) {
	response := RawResponse{
		StatusCode: 200,
		Body:       []byte(MOCK_BODY),
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Code != 200 {
		t.Errorf(EXPECTED_DIGIT_ERROR, 200, writer.Code)
	}
}

func TestObjectResponseWritesStatusCode(t *testing.T) {
	response := ObjectResponse{
		StatusCode: 200,
		Body:       map[string]interface{}{},
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Code != 200 {
		t.Errorf(EXPECTED_DIGIT_ERROR, 200, writer.Code)
	}
}

func TestObjectResponseSetsJSONResponseType(t *testing.T) {
	response := ObjectResponse{
		StatusCode: 200,
		Body:       map[string]interface{}{},
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Header().Get(CONTENT_TYPE_HEADER_KEY) != CONTENT_TYPE_JSON {
		t.Errorf(EXPECTED_STRING_ERROR, CONTENT_TYPE_JSON, writer.Header().Get(CONTENT_TYPE_HEADER_KEY))
	}
}

func TestObjectResponseJSONEncodesObject(t *testing.T) {
	response := ObjectResponse{
		StatusCode: 200,
		Body:       map[string]string{"#": MOCK_BODY},
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Body.String() != MOCK_BODY_JSON {
		t.Errorf(EXPECTED_STRING_ERROR, MOCK_BODY_JSON, writer.Body.String())
	}
}

func TestObjectResponseRaisesInternalServerErrorOnJSONEncodingError(t *testing.T) {
	response := ObjectResponse{
		StatusCode: 200,
		Body:       make(chan int),
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Code != 500 {
		t.Errorf(EXPECTED_DIGIT_ERROR, 500, writer.Code)
	}
	bytes := writer.Body.Bytes()
	expected := "internal server error: json: unsupported type: chan int"
	if string(bytes) != expected {
		t.Errorf(EXPECTED_STRING_ERROR, expected, writer.Body.String())
	}
}

func TestJsonResponseWritesStatusCode(t *testing.T) {
	response := JSONResponse{
		StatusCode: 200,
		Body:       map[string]interface{}{},
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Code != 200 {
		t.Errorf(EXPECTED_DIGIT_ERROR, 200, writer.Code)
	}
}

func TestJsonResponseSetsJSONResponseType(t *testing.T) {
	response := JSONResponse{
		StatusCode: 200,
		Body:       map[string]interface{}{},
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Header().Get(CONTENT_TYPE_HEADER_KEY) != CONTENT_TYPE_JSON {
		t.Errorf("Expected %s, got %s", CONTENT_TYPE_JSON, writer.Header().Get(CONTENT_TYPE_HEADER_KEY))
	}
}

func TestJsonResponseJSONEncodesObject(t *testing.T) {
	response := JSONResponse{
		StatusCode: 200,
		Body:       map[string]interface{}{"#": MOCK_BODY},
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Body.String() != MOCK_BODY_JSON {
		t.Errorf(EXPECTED_STRING_ERROR, MOCK_BODY_JSON, writer.Body.String())
	}
}

func TestErrorResponseWritesStatusCode(t *testing.T) {
	response := ErrorResponse{
		StatusCode: 503,
		Message:    "service unavailable",
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Code != 503 {
		t.Errorf(EXPECTED_DIGIT_ERROR, 503, writer.Code)
	}
}

func TestErrorResponseSetsMessage(t *testing.T) {
	response := ErrorResponse{
		StatusCode: 503,
		Message:    SERVICE_UNAVAILABLE,
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Body.String() != SERVICE_UNAVAILABLE {
		t.Errorf("Expected 'service unavailable', got %s", writer.Body.String())
	}
}

func TestErrorResponseSetsMessageAndError(t *testing.T) {
	response := ErrorResponse{
		StatusCode: 503,
		Message:    SERVICE_UNAVAILABLE,
		Error:      fmt.Errorf("error"),
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Body.String() != SERVICE_UNAVAILABLE+": error" {
		t.Errorf(EXPECTED_STRING_ERROR, SERVICE_UNAVAILABLE, writer.Body.String())
	}
}

func TestBadRequestResponseWritesStatusCode(t *testing.T) {
	response := BadRequest{}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Code != 400 {
		t.Errorf(EXPECTED_DIGIT_ERROR, 400, writer.Code)
	}
}

func TestBadRequestSetsMessage(t *testing.T) {
	response := BadRequest{}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Body.String() != "bad request" {
		t.Errorf("Expected 'bad request', got %s", writer.Body.String())
	}
}

func TestBadRequestSetsMessageAndError(t *testing.T) {
	response := BadRequest{
		Error: fmt.Errorf("error"),
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Body.String() != "bad request: error" {
		t.Errorf("Expected 'bad request: error', got %s", writer.Body.String())
	}
}

func TestInternalServerErrorResponseWritesStatusCode(t *testing.T) {
	response := InternalServerError{}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Code != 500 {
		t.Errorf(EXPECTED_DIGIT_ERROR, 500, writer.Code)
	}
}

func TestInternalServerErrorSetsMessage(t *testing.T) {
	response := InternalServerError{}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Body.String() != "internal server error" {
		t.Errorf(EXPECTED_STRING_ERROR, INTERNAL_SERVER_ERROR, writer.Body.String())
	}
}

func TestInternalServerErrorSetsMessageAndError(t *testing.T) {
	response := InternalServerError{
		Error: fmt.Errorf("error"),
	}
	writer := httptest.NewRecorder()
	response.Write(writer)
	expected := "internal server error: error"
	if writer.Body.String() != expected {
		t.Errorf(EXPECTED_STRING_ERROR, expected, writer.Body.String())
	}
}

func TestServiceUnavailableSetsStatusCode(t *testing.T) {
	response := ServiceUnavailable{}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Code != 503 {
		t.Errorf(EXPECTED_DIGIT_ERROR, 503, writer.Code)
	}
}

func TestServiceNotAvailableSetsMessage(t *testing.T) {
	response := ServiceUnavailable{}
	writer := httptest.NewRecorder()
	response.Write(writer)
	if writer.Body.String() != SERVICE_UNAVAILABLE {
		t.Errorf(EXPECTED_STRING_ERROR, SERVICE_UNAVAILABLE, writer.Body.String())
	}
}
