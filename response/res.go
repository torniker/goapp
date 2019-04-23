package response

import (
	"github.com/torniker/wrap/logger"
)

// NewResponse create an instance of sub responder
func NewResponse() *Res {
	return &Res{
		headers: make(map[string]string),
	}
}

// Res is type for sub responses
type Res struct {
	headers   map[string]string
	status    int
	committed bool
	output    interface{}
}

// SetStatus sets response status code
func (r *Res) SetStatus(status int) {
	r.status = status
}

// SetHeader sets header for responder
func (r *Res) SetHeader(key, val string) {
	r.headers[key] = val
}

// Status returns response status code
func (r *Res) Status() int {
	return r.status
}

// Commited returns response status
func (r *Res) Commited() bool {
	return r.committed
}

// Write commits and writes data into the response body
func (r *Res) Write(body interface{}) error {
	if r.status == 0 {
		r.SetStatus(200)
	}
	r.committed = true
	r.output = body
	logger.Infof("<--- status: %v, body: %#v", r.status, body)
	return nil
	// return json.NewEncoder(os.Stdout).Encode(body)
}

// Output returns response output
func (r *Res) Output() interface{} {
	return r.output
}

// EnableCORS sets corresponsing headers to enable CORS
func (r *Res) EnableCORS(origin, methods, headers string) {
	r.SetHeader("Access-Control-Allow-Origin", origin)
	r.SetHeader("Access-Control-Allow-Methods", methods)
	r.SetHeader("Access-Control-Allow-Headers", headers)
}
