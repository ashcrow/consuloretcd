package consuloretcd

// Error map
// The key is the error that is returned from a failed call
// The value is a string description
var Errors map[int]string = map[int]string{
	1: "Unable to contact remote server",
	2: "Unexpected server status code in response",
	3: "Unable to read response body",
	4: "Unable to decode the value response",
	5: "Server did not save the new key",
}

// Interface to be a valid KeyValueStore
type KeyValueStore interface {
	makeUri() string
	GetKey() (KeyValue, interface{})
	PutKey() (KeyValue, interface{})
}

// Key/Value abstraction used
type KeyValue struct {
	Name        string
	Key         string
	Exists      bool
	Error       int
	StatusCode  int
	CreateIndex int
	ModifyIndex int
	Value       interface{}
}
