package consuloretcd

import (
	"fmt"
	"reflect"
	"testing"
)

func ExampleNewClient_consul() {
	// error checking omitted
	consul, _ := NewClient(
		"consul",
		ConsulDefaultConfig,
	)
	fmt.Println(reflect.TypeOf(consul))
	// Output: consuloretcd.Consul
}

func ExampleNewClient_etcd() {
	// error checking omitted
	etcd, _ := NewClient(
		"etcd",
		EtcdDefaultConfig,
	)
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
		EtcdDefaultConfig,
	)
	keyval, _ := etcd.GetKey("keyname", KeyOptions{})

	fmt.Println(reflect.TypeOf(keyval))
	// Output: consuloretcd.KeyValue
}

func ExampleEtcd_PutKey() {
	// error checking omitted
	etcd, _ := NewClient(
		"etcd",
		EtcdDefaultConfig,
	)
	keyval, _ := etcd.PutKey("keyname", "a value", KeyOptions{})

	fmt.Println(reflect.TypeOf(keyval))
	// Output: consuloretcd.KeyValue
}

func ExampleEtcd_DeleteKey() {
	// error checking omitted
	etcd, _ := NewClient(
		"etcd",
		EtcdDefaultConfig,
	)

	etcd.DeleteKey("keyname", KeyOptions{})
}

func ExampleConsul_GetKey() {
	// error checking omitted
	consul, _ := NewClient(
		"consul",
		ConsulDefaultConfig,
	)

	keyval, _ := consul.GetKey("keyname", KeyOptions{})

	fmt.Println(reflect.TypeOf(keyval))
	// Output: consuloretcd.KeyValue
}

func ExampleConsul_PutKey() {
	// error checking omitted
	consul, _ := NewClient(
		"consul",
		ConsulDefaultConfig,
	)

	keyval, _ := consul.PutKey("keyname", "a value", KeyOptions{})

	fmt.Println(reflect.TypeOf(keyval))
	// Output: consuloretcd.KeyValue
}

func ExampleConsul_DeleteKey() {
	// error checking omitted
	consul, _ := NewClient(
		"consul",
		ConsulDefaultConfig,
	)

	consul.DeleteKey("keyname", KeyOptions{})
}
