package web

import "net/http"

// ProxyResponseWriter records the status code written by a call to the
// WriteHeader function on a http.ResponseWriter interface. This type also
// implements the http.ResponseWriter interface.
type ProxyResponseWriter struct {
	Status          int
	UpstreamHeaders http.Header
	http.ResponseWriter
}

// Header implements the http.ResponseWriter interface and simply relays the
// request.
func (prw *ProxyResponseWriter) Header() http.Header {
	return prw.ResponseWriter.Header()
}

// Write implements the http.ResponseWriter interface and simply relays the
// request.
func (prw *ProxyResponseWriter) Write(data []byte) (int, error) {
	return prw.ResponseWriter.Write(data)
}

// WriteHeader implements the http.ResponseWriter interface and simply relays
// the request after cleaning up the request headers. It theb records the status
// code written.
func (prw *ProxyResponseWriter) WriteHeader(status int) {

	// After the header is written, the headers will be flushed, so at this point
	// we have to strip all the headers which should not be written upstream.
	for k, vv := range prw.UpstreamHeaders {

		// Remove the headers from the request that match the upstream headers.
		prw.Header().Del(k)

		for _, v := range vv {

			// Add in the header from the upstream request to fill in what we had
			// before we proxied.
			prw.Header().Add(k, v)
		}
	}

	prw.ResponseWriter.WriteHeader(status)
	prw.Status = status
}
