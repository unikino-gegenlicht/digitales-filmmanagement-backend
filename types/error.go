package types

import (
	"net/http"
)

// APIError is the default response layout for an error that occurred during
// request handling.
type APIError struct {
	// ErrorCode contains a short error code identifying the error and the
	// service
	ErrorCode string `json:"code"`
	// ErrorTitle contains a human-readable error title
	ErrorTitle string `json:"title"`
	// ErrorDescription contains a human-readable description of the error that
	// occurred while handling the request
	ErrorDescription string `json:"description"`
	// HttpStatusCode contains the numeric http code that is associated with the
	// error. this value should be used to write out the error
	HttpStatusCode int `json:"httpCode"`
	// HttpStatusText contains the description of the HttpStatusCode to allow
	// humans to understand the numeric http code
	HttpStatusText string `json:"httpError"`
}

// InferHttpStatusText takes the already configured HTTP Status Code and infers
// the HTTP status text from it using the net/http package. If the status code
// is not known to the package the HTTP status text will stay empty
func (e *APIError) InferHttpStatusText() {
	e.HttpStatusText = http.StatusText(e.HttpStatusCode)
}

// WrapError takes a native golang error as parameter and wraps it into a
// WISdoMError. The WISdoMError instance will overwrite every field already
// present on the error and set the status code to 500 to indicate a internal
// error occurred. It optionally takes a service name as argument, if multiple
// names are supplied they are joined together using a dot (.)
func (e *APIError) WrapError(err error) {
	// now build full the full error code used for wrapping internal errors
	e.ErrorCode = "INTERNAL_ERROR"
	// set the title to a generic error title
	e.ErrorTitle = "Internal Error in Backend"
	// set the error description to the external error
	e.ErrorDescription = err.Error()
	// set the http code to 500 and infer the text from this
	e.HttpStatusCode = http.StatusInternalServerError
	e.HttpStatusText = http.StatusText(e.HttpStatusCode)
}
