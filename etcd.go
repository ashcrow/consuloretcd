package consuloretcd

// Etcd specific implementation.

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// EtcdDefaultconfig defines sane defaults for Etcd
var EtcdDefaultConfig Config = Config{
	Endpoint: "http://127.0.0.1",
	Client:   http.Client{},
	Port:     4001,
}

// Etcd is the useable structure for Etcd
type Etcd struct {
	Config
}

// makeURI makes the URI from the Etcd struct
// Returns the full URI as a string
func (c Etcd) makeURI(name string, opts KeyOptions) string {
	return c.Endpoint + ":" + strconv.Itoa(c.Port) + "/v2/keys/" + name + "?"
}

// makeParams creates a url.Values instance based off the KeyOptions
func (c Etcd) makeParams(opts KeyOptions) url.Values {
	v := url.Values{}
	// TODO(ashcrow): This is a hack to avoid colliding with int:0. Fix it.
	if opts.CASet != "" {
		v.Set("prevIndex", opts.CASet)
	}
	if opts.TTL != 0 {
		v.Set("ttl", strconv.Itoa(opts.TTL))
	}
	return v
}

// GetKey gets a key from the remote Etcd server.
// Returns KeyValue, nil on success
// Returns KeyValue, int (lookup via Errors) when unable to get a value
func (c Etcd) GetKey(name string, opts KeyOptions) (KeyValue, error) {
	kv := KeyValue{
		Name:   name,
		Exists: false}
	params := c.makeParams(opts)
	resp, err := c.Client.Get(c.makeURI(name, opts) + "?" + params.Encode())
	if err != nil {
		kv.Error = 1
		return kv, errors.New(Errors[kv.Error])
	}
	// Close the body at the end
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		kv.Error = 2
		kv.StatusCode = resp.StatusCode
		return kv, errors.New(Errors[kv.Error])
	}
	kv.StatusCode = resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		kv.Error = 3
		return kv, errors.New(Errors[kv.Error])
	}
	// Unmarshal the response
	var result interface{}
	json.Unmarshal(body, &result)

	nr := result.(map[string]interface{})
	node := nr["node"].(map[string]interface{})

	// Turn the result into a map of interfaces
	kv.CreateIndex = int(node["createdIndex"].(float64))
	kv.ModifyIndex = int(node["modifiedIndex"].(float64))
	kv.Key = node["key"].(string)
	kv.Error = 0
	kv.Exists = true
	kv.Value = node["value"]
	return kv, nil
}

// PutKey puts a key on the remote Etcd server.
// Returns KeyValue, nil on success
// Returns KeyValue, ERROR_CODE when unable to get a value
func (c Etcd) PutKey(name string, value string, opts KeyOptions) (KeyValue, error) {
	kv := KeyValue{
		Name:   name,
		Exists: false}
	params := c.makeParams(opts)
	params.Add("value", value)

	req, _ := http.NewRequest(
		"PUT",
		c.makeURI(name, opts),
		strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded; param=value")

	resp, err := c.Client.Do(req)
	if err != nil {
		kv.Error = 1
		return kv, errors.New(Errors[kv.Error])
	}
	// Close the body at the end
	defer resp.Body.Close()
	kv.StatusCode = resp.StatusCode
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		kv.Exists = false
		kv.Error = 2
		return kv, errors.New(Errors[kv.Error])
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		kv.Error = 3
		return kv, errors.New(Errors[kv.Error])
	}

	// Unmarshal the response
	var result interface{}
	json.Unmarshal(body, &result)

	nr := result.(map[string]interface{})

	if nr["action"] == "set" {
		node := nr["node"].(map[string]interface{})
		kv.Value = node["value"]
		kv.Exists = true
		kv.Key = node["key"].(string)
		kv.ModifyIndex = int(node["modifiedIndex"].(float64))
		kv.CreateIndex = int(node["createdIndex"].(float64))
		return kv, nil
	}
	kv.Error = 5
	return kv, errors.New(Errors[kv.Error])
}

// DeleteKey deletes a key from the remote Etcd server.
// Returns nil on success
// Returns Error when unable to delete
func (c Etcd) DeleteKey(name string, opts KeyOptions) error {
	req, _ := http.NewRequest("DELETE", c.makeURI(name, opts), nil)
	resp, err := c.Client.Do(req)
	if err != nil {
		return errors.New(Errors[6])
	}
	// Close the body at the end
	defer resp.Body.Close()
	// It's weird but 200 means it deleted ...
	if resp.StatusCode != 200 {
		return errors.New(Errors[6])
	}
	return nil
}
