package wapp

// Enums

type DataType string
const (
	DataTypeHTML DataType = "text/html"
	DataTypeJSON DataType = "application/json"
	DataTypeXML  DataType = "text/xml"
)

type HTTPMethod string
const (
	HTTPMethodAll     HTTPMethod = "ALL"
	HTTPMethodGet     HTTPMethod = "GET"
	HTTPMethodPost    HTTPMethod = "POST"
	HTTPMethodPut     HTTPMethod = "PUT"
	HTTPMethodGDelete HTTPMethod = "DELETE"
	HTTPMethodHead    HTTPMethod = "HEAD"
	HTTPMethodConnect HTTPMethod = "CONNECT"
	HTTPMethodOptions HTTPMethod = "OPTIONS"
	HTTPMethodTrace   HTTPMethod = "TRACE"
)