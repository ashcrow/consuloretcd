/*
A simplistic key/value abstraction for use with Consul and Etcd.
*/
package consuloretcd

// Items shared between implementations
import (
	"net/http"
)

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
	makeURI(string) string
	GetKey(string) (KeyValue, error)
	PutKey(string, string) (KeyValue, error)
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

// Configuration for a KeyValueStore
type Config struct {
	Endpoint string
	Client   http.Client
	Port     int
}

// Returns a new KeyValueStore client based on the name
func NewClient(name string, config Config) (KeyValueStore, interface{}) {
	switch {
	case name == "consul":
		return Consul{
			KeyValueStoreConfig: config,
		}, nil
	case name == "etcd":
		return Etcd{
			KeyValueStoreConfig: config,
		}, nil
	default:
		return nil, 1
	}
}
