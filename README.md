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

## Documentation
Read the docs at godoc: http://godoc.org/github.com/ashcrow/consuloretcd

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

    // Consul example. Replace "consul" with "etcd" to use etcd.
    consul, _ := consuloretcd.NewClient(
        "consul",
		consuloretcd.ConsulDefaultConfig)

    // Get a key in consul
	consul_res1, _ := consul.GetKey("test")
	fmt.Println(consul_res1)

    // Set a key in consul
	consul_res2,  _ := consul.PutKey("test", "saa")
	fmt.Println(consul_res2)

    // Delete a key in consul
    if err := consul.DeleteKey("test"); err != nil {
	    fmt.Println(err)
    }
}
