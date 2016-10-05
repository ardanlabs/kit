// Package web provides web application support for context and MongoDB access.
// Current Status Codes:
//		200 OK           : StatusOK                  : Call is success and returning data.
//		204 No Content   : StatusNoContent           : Call is success and returns no data.
//		400 Bad Request  : StatusBadRequest          : Invalid post data (syntax or semantics).
//		401 Unauthorized : StatusUnauthorized        : Authentication failure.
//		404 Not Found    : StatusNotFound            : Invalid URL or identifier.
//		500 Internal     : StatusInternalServerError : Weblication specific beyond scope of user.
package web

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// Invalid describes a validation error belonging to a specific field.
type Invalid struct {
	Fld string `json:"field_name"`
	Err string `json:"error"`
}

// jsonError is the response for errors that occur within the API.
type jsonError struct {
	Error  string    `json:"error"`
	Fields []Invalid `json:"fields,omitempty"`
}

//==============================================================================

// Context contains data associated with a single request.
type Context struct {
	http.ResponseWriter
	Request   *http.Request
	Now       time.Time
	Params    map[string]string
	SessionID string
	Status    int
	Ctx       map[string]interface{}
	Web       *Web
}

// Error handles all error responses for the API.
func (c *Context) Error(err error) {
	switch err {
	case ErrNotFound:
		c.RespondError(err.Error(), http.StatusNotFound)
	case ErrInvalidID:
		c.RespondError(err.Error(), http.StatusBadRequest)
	case ErrValidation:
		c.RespondError(err.Error(), http.StatusBadRequest)
	case ErrNotAuthorized:
		c.RespondError(err.Error(), http.StatusUnauthorized)
	default:
		c.RespondError(err.Error(), http.StatusInternalServerError)
	}
}

// Respond sends JSON to the client.
// If code is StatusNoContent, v is expected to be nil.
func (c *Context) Respond(data interface{}, code int) error {
	c.Status = code

	// Just set the status code and we are done.
	if code == http.StatusNoContent {
		c.WriteHeader(code)
		return nil
	}

	// Set the content type.
	c.Header().Set("Content-Type", "application/json")

	// Write the status code.
	c.WriteHeader(code)

	// Marshal the data into a JSON string.
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Look for a JSONP marker
	if cb := c.Request.URL.Query().Get("callback"); cb != "" {

		// We need to wrap the result in a function call.
		// callback_value({"data_1": "hello world", "data_2": ["the","sun","is","shining"]});
		io.WriteString(c, cb+"("+string(jsonData)+")")
		return nil
	}

	// We can send the result straight through.
	io.WriteString(c, string(jsonData))

	return nil
}

// RespondInvalid sends JSON describing field validation errors.
func (c *Context) RespondInvalid(fields []Invalid) {
	v := jsonError{
		Error:  "field validation failure",
		Fields: fields,
	}
	c.Respond(v, http.StatusBadRequest)
}

// RespondError sends JSON describing the error
func (c *Context) RespondError(error string, code int) {
	c.Respond(jsonError{Error: error}, code)
}

// Proxy will setup a direct proxy inbetween this service and the destination
// service.
func (c *Context) Proxy(targetURL string, rewrite func(req *http.Request)) error {
	target, err := url.Parse(targetURL)
	if err != nil {
		return err
	}

	// Define our custom request director to ensure that the correct headers are
	// forwarded as well as having the request path and query rewritten
	// properly.
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)

		if target.RawQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = target.RawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = target.RawQuery + "&" + req.URL.RawQuery
		}

		// Set the headers on the incoming request by clearing them all out.
		req.Header = make(http.Header)

		// Rewrite the request for the director if a rewrite function was passed.
		if rewrite != nil {
			rewrite(req)
		}
	}

	// Create a new reverse proxy. We need to to this here because for the path
	// rewriting we may need access to variables stored in this specific
	// request's path parameters which can allow to be overridden via the
	// rewrite argument to this function.
	proxy := httputil.ReverseProxy{Director: director}

	// Create a new proxy response writer that will record the http status code
	// issued by the reverse proxy.
	prw := ProxyResponseWriter{
		ResponseWriter: c.ResponseWriter,
	}

	// Serve the request via the built in handler here.
	proxy.ServeHTTP(&prw, c.Request)

	c.Status = prw.Status

	return nil
}

// singleJoiningSlash ensures that there is a single joining slash inbetween the
// url's that are being joined. This was sourced from the reverseproxy.go file
// inside the net/http/httputil package in the stdlib.
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")

	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}

	return a + b
}
