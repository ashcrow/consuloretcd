package consuloretcd

func ExampleNewClient_consul() {
	consul, err := consuloretcd.NewClient(
		"consul",
		"http://127.0.0.1",
		http.Client{},
		8500)
}

func ExampleNewClient_etcd() {
	etcd, err := consuloretcd.NewClient(
		"etcd",
		"http://127.0.0.1",
		client,
		4001)
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
