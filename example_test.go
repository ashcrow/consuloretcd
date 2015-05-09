package consuloretcd

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func ExampleNewClient_consul() {
	// error checking omitted
	consul, _ := NewClient(
		"consul",
		Config{
			Endpoint: "http://127.0.0.1",
			Client:   http.Client{},
			Port:     8500})
	fmt.Println(reflect.TypeOf(consul))
	// Output: consuloretcd.Consul
}

func ExampleNewClient_etcd() {
	// error checking omitted
	etcd, _ := NewClient(
		"etcd",
		Config{
			Endpoint: "http://127.0.0.1",
			Client:   http.Client{},
			Port:     4001})
	fmt.Println(reflect.TypeOf(etcd))
	// Output: consuloretcd.Etcd
}

func TestNewClient_default(t *testing.T) {
	unknown, err := NewClient(
		"unknown",
		Config{})
	if unknown != nil || err == nil {
		t.FailNow()
	}
}

func ExampleEtcd_GetKey() {
	// error checking omitted
	etcd, _ := NewClient(
		"etcd",
		Config{
			Endpoint: "http://127.0.0.1",
			Client:   http.Client{},
			Port:     4001})

	keyval, _ := etcd.GetKey("keyname", KeyOptions{})

	fmt.Println(reflect.TypeOf(keyval))
	// Output: consuloretcd.KeyValue
}

func ExampleEtcd_PutKey() {
	// error checking omitted
	etcd, _ := NewClient(
		"etcd",
		Config{
			Endpoint: "http://127.0.0.1",
			Client:   http.Client{},
			Port:     4001})

	keyval, _ := etcd.PutKey("keyname", "a value", KeyOptions{})

	fmt.Println(reflect.TypeOf(keyval))
	// Output: consuloretcd.KeyValue
}

func ExampleEtcd_DeleteKey() {
	// error checking omitted
	etcd, _ := NewClient(
		"etcd",
		Config{
			Endpoint: "http://127.0.0.1",
			Client:   http.Client{},
			Port:     4001})

	etcd.DeleteKey("keyname", KeyOptions{})
}

func ExampleConsul_GetKey() {
	// error checking omitted
	consul, _ := NewClient(
		"consul",
		Config{
			Endpoint: "http://127.0.0.1",
			Client:   http.Client{},
			Port:     8500})

	keyval, _ := consul.GetKey("keyname", KeyOptions{})

	fmt.Println(reflect.TypeOf(keyval))
	// Output: consuloretcd.KeyValue
}

func ExampleConsul_PutKey() {
	// error checking omitted
	consul, _ := NewClient(
		"consul",
		Config{
			Endpoint: "http://127.0.0.1",
			Client:   http.Client{},
			Port:     8500})

	keyval, _ := consul.PutKey("keyname", "a value", KeyOptions{})

	fmt.Println(reflect.TypeOf(keyval))
	// Output: consuloretcd.KeyValue
}

func ExampleConsul_DeleteKey() {
	// error checking omitted
	consul, _ := NewClient(
		"consul",
		Config{
			Endpoint: "http://127.0.0.1",
			Client:   http.Client{},
			Port:     8500})

	consul.DeleteKey("keyname", KeyOptions{})
}
