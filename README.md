# consuloretcd

A simplistic key/value abstraction for use with Consul and Etcd.

**Warning**: This needs a lot of cleaning up and code consolidation.

**Warning**: The api will probably change a lot.

## Install

```bash
go get github.com/ashcrow/consuloretcd
```

## License
See LICENSE

## Example
```go
package main

import (
	"fmt"
	"github.com/ashcrow/consuloretcd"
	"net/http"
)

// Example
func main() {
    // You must provide and http.Client
	client := http.Client{}

    // Consul example
    consul, _ := consuloretcd.NewClient(
        "consul",
        "http://127.0.0.1",
        client,
        8500)

    // Get a key in consul
	consul_res1, _ := consul.GetKey("test")
	fmt.Println(consul_res1)
    // Set a key in consul
	consul_res2,  _ := consul.PutKey("test", "saa")
	fmt.Println(consul_res2)

    // Etcd example
    etcd, _ := consuloretcd.NewClient(
        "etcd",
        "http://127.0.0.1",
        client,
        4001)

    // Get a key from etcd
	etcd_res1, _ := etcd.GetKey("test")
	fmt.Println(etcd_res1)
    // Set a key in etcd
	etcd_res2, _ := etcd.PutKey("test", "saa")
	fmt.Println(etcd_res2)
}
```
