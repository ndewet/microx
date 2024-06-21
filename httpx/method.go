package httpx

type Method string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	DELETE  Method = "DELETE"
	PATCH   Method = "PATCH"
	PUT     Method = "PUT"
	HEAD    Method = "HEAD"
	OPTIONS Method = "OPTIONS"
	CONNECT Method = "CONNECT"
	TRACE   Method = "TRACE"
)
