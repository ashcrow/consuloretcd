package consuloretcd

func ExampleNewClient_consul() {
	consul, err := consuloretcd.NewClient(
		"consul",
		consuloretcd.Config{
			Endpoint: "http://127.0.0.1",
			Client:   http.Client{},
			Port:     8500})
}

func ExampleNewClient_etcd() {
	etcd, err := consuloretcd.NewClient(
		"etcd",
		consuloretcd.Config{
			Endpoint: "http://127.0.0.1",
			Client:   http.Client{},
			Port:     4001})
}

func ExampleEtcd_GetKey() {
	keyval, err := consul.GetKey("keyname")
}

func ExampleEtcd_PutKey() {
	keyval, err := etcd.PutKey("keyname", "a value")
}

func ExampleConsul_GetKey() {
	keyval, err := consul.GetKey("keyname")
}

func ExampleConsul_PutKey() {
	keyval, err := consul.PutKey("keyname", "a value")
}
