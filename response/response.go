package response

// list of response types
const (
	TypeHTTP int = 1 + iota
	TypeCLI
)

// Response defines methods for response object
type Response interface {
	SetStatus(status int)
	Status() int
	SetHeader(key, val string)
	Commited() bool
	Write(interface{}) error
	SuccessWithNoContent()
	Output() interface{}
	EnableCORS(origin, methods, headers string)
}
